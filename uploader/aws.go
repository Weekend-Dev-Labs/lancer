package uploader

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/weekend-dev-labs/lancer/types"
)

type AwsUploader struct {
	mu       sync.Mutex
	s3Client *s3.Client
	sessions map[string]MultipartSessionInfo
}

type MultipartSessionInfo struct {
	UploadInfo *s3.CreateMultipartUploadOutput
	Parts      s3Types.CompletedPart
}

func NewAwsUploader() *AwsUploader {
	config, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithSharedConfigFiles([]string{config.path}),
	)

	client := s3.NewFromConfig(config)

	return client
}

func (au *AwsUploader) CreateMultipart(bucket string, key string) (*s3.CreateMultipartUploadOutput, error) {
	return au.s3Client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket: &bucket,
		Key:    &key,
	})
}

type UploadPartParam struct {
	Bucket   string
	Key      string
	Part     int32
	Body     *bytes.Reader
	UploadID string
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

type CompleteMultipartParam struct {
	Bucket             string
	Key                string
	UploadID           string
	CompletedPartsInfo []s3Types.CompletedPart
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

type AbortMultipartParam struct {
	Bucket   string
	Key      string
	UploadID string
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

}
