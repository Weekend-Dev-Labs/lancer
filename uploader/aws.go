package uploader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/types"
)

type AwsUploader struct {
	mu       sync.Mutex
	s3Client *s3.Client
	sessions map[string]*MultipartSessionInfo
}

type MultipartSessionInfo struct {
	UploadInfo *s3.CreateMultipartUploadOutput
	Parts      s3Types.CompletedMultipartUpload
}

type AbortMultipartParam struct {
	Bucket   string
	Key      string
	UploadID string
}

type UploadPartParam struct {
	Bucket   string
	Key      string
	Part     int32
	Body     *bytes.Reader
	UploadID string
}

type CompleteMultipartParam struct {
	Bucket             string
	Key                string
	UploadID           string
	CompletedPartsInfo []s3Types.CompletedPart
}

type UploadFullFileParam struct {
	Bucket string
	Key    string
	File   io.Reader
}

func NewAwsUploader(cfg *config.LancerConfig) *AwsUploader {
	if cfg.Store.AWS.Store {
		config, err := awsConfig.LoadDefaultConfig(
			context.TODO(),
			awsConfig.WithSharedConfigFiles([]string{cfg.Store.AWS.Config}),
			awsConfig.WithRegion(cfg.Store.AWS.Region),
		)

		if err != nil {
			log.Fatalf("[LANCER ERROR] Invalid AWS Configuration (%v)", err.Error())
		}

		client := s3.NewFromConfig(config)

		return &AwsUploader{
			s3Client: client,
			mu:       sync.Mutex{},
			sessions: make(map[string]*MultipartSessionInfo),
		}
	}

	return nil
}

func (au *AwsUploader) CreateMultipart(bucket string, key string, sessionInfo *types.SessionInfo) error {
	uploadOutput, err := au.s3Client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		return err
	}

	au.mu.Lock()
	au.sessions[sessionInfo.ID] = &MultipartSessionInfo{
		UploadInfo: uploadOutput,
		Parts:      s3Types.CompletedMultipartUpload{},
	}
	au.mu.Unlock()

	return nil
}

func (au *AwsUploader) UploadPart(params UploadPartParam) (*s3.UploadPartOutput, error) {
	input := &s3.UploadPartInput{
		Bucket:     aws.String(params.Bucket),
		Key:        aws.String(params.Key),
		PartNumber: &params.Part,
		UploadId:   &params.UploadID,
		Body:       params.Body,
	}

	resp, err := au.s3Client.UploadPart(context.TODO(), input)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (au *AwsUploader) CompleteMultipartUpload(params CompleteMultipartParam) error {
	input := &s3.CompleteMultipartUploadInput{
		Bucket:   &params.Bucket,
		Key:      &params.Key,
		UploadId: &params.UploadID,
		MultipartUpload: &s3Types.CompletedMultipartUpload{
			Parts: params.CompletedPartsInfo,
		},
	}

	_, err := au.s3Client.CompleteMultipartUpload(context.TODO(), input)

	return err
}

func (au *AwsUploader) AbortMultipartUpload(params AbortMultipartParam) error {
	input := &s3.AbortMultipartUploadInput{
		Bucket:   &params.Bucket,
		Key:      &params.Key,
		UploadId: &params.UploadID,
	}

	_, err := au.s3Client.AbortMultipartUpload(context.TODO(), input)

	return err
}

func (au *AwsUploader) HandleMultipartUploads(id string, part int32, sessionInfo *types.SessionInfo, file *bytes.Reader) error {
	au.mu.Lock()
	info, isExists := au.sessions[id]
	au.mu.Unlock()

	if !isExists {
		return fmt.Errorf("no upload session found")
	}

	if sessionInfo.CurrentChunk+1 == int64(part) {
		err := au.CompleteMultipartUpload(CompleteMultipartParam{
			Bucket:             *info.UploadInfo.Bucket,
			Key:                *info.UploadInfo.Key,
			UploadID:           *info.UploadInfo.UploadId,
			CompletedPartsInfo: info.Parts.Parts,
		})

		if err != nil {
			return err
		}
	}

	upload, err := au.UploadPart(UploadPartParam{
		Bucket:   *info.UploadInfo.Bucket,
		Key:      *info.UploadInfo.Key,
		Part:     part,
		Body:     file,
		UploadID: *info.UploadInfo.UploadId,
	})

	if err != nil {
		return err
	}

	info.Parts.Parts = append(info.Parts.Parts, s3Types.CompletedPart{
		ETag:       upload.ETag,
		PartNumber: &part,
	})

	au.mu.Lock()
	au.sessions[id] = info
	au.mu.Unlock()

	return nil
}

func (au *AwsUploader) UploadFullFile(params *UploadFullFileParam) (*s3.PutObjectOutput, error) {
	return au.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &params.Bucket,
		Key:    &params.Key,
		Body:   params.File,
	})
}
