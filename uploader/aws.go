package uploader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/utils"
)

type AwsUploader struct {
	mu       sync.Mutex
	s3Client *s3.Client
	config   *config.LancerConfig
	sessions map[string]*MultipartSessionInfo
}

type MultipartSessionInfo struct {
	UploadInfo *s3.CreateMultipartUploadOutput
	Parts      s3Types.CompletedMultipartUpload
}

type AbortMultipartParam struct {
	Bucket    string
	Key       string
	UploadID  string
	SessionID string
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

type Metadata struct {
	Key    string `json:"key"`
	Etag   string `json:"etag"`
	Bucket string `json:"bucket"`
}

func NewAwsUploader(cfg *config.LancerConfig) *AwsUploader {
	if cfg.Store.AWS.Store {
		config, err := awsConfig.LoadDefaultConfig(
			context.TODO(),
			awsConfig.WithRegion(cfg.Store.AWS.Region),
			awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.Store.AWS.ClientID, cfg.Store.AWS.ClientSecret, "")),
		)

		if err != nil {
			log.Fatalf("[LANCER ERROR] Invalid AWS Configuration (%v)", err.Error())
		}

		client := s3.NewFromConfig(config)

		return &AwsUploader{
			s3Client: client,
			mu:       sync.Mutex{},
			sessions: make(map[string]*MultipartSessionInfo),
			config:   cfg,
		}
	}

	return nil
}

func (au *AwsUploader) CreateChunkUploadSession(sessionInfo *types.SessionInfo) error {
	fileKey := fmt.Sprintf("%d_%s", time.Now().Unix(), sessionInfo.FileName)

	upload, err := au.s3Client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket: &au.config.Store.AWS.Bucket,
		Key:    &fileKey,
	})

	if err != nil {
		return err
	}

	au.mu.Lock()
	au.sessions[sessionInfo.ID] = &MultipartSessionInfo{
		UploadInfo: upload,
		Parts:      s3Types.CompletedMultipartUpload{},
	}
	au.mu.Unlock()

	return nil
}

func (au *AwsUploader) Upload(sessionInfo *types.SessionInfo, file []byte) (*UploadAck, error) {
	resp, err := au.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &au.config.Store.AWS.Bucket,
		Key:    &sessionInfo.FileName,
		Body:   bytes.NewReader(file),
	})

	if err != nil {
		return nil, err
	}

	return &UploadAck{
		Provider: types.UploaderAws,
		ProviderMetadata: map[string]string{
			"bucket": au.config.Store.AWS.Bucket,
			"key":    sessionInfo.FileName,
			"etag":   *resp.ETag,
		},
		Checksum: "",
		FilePath: "",
	}, nil
}

func (au *AwsUploader) HandlePartUpload(sessionInfo *types.SessionInfo, file []byte) error {
	au.mu.Lock()
	info, isExists := au.sessions[sessionInfo.ID]
	au.mu.Unlock()

	if !isExists {
		return fmt.Errorf("no upload session found")
	}

	partNumber := int32(sessionInfo.CurrentChunk)
	partNumber += 1

	resp, err := au.s3Client.UploadPart(context.TODO(), &s3.UploadPartInput{
		Bucket:     info.UploadInfo.Bucket,
		Key:        info.UploadInfo.Key,
		UploadId:   info.UploadInfo.UploadId,
		PartNumber: &partNumber,
		Body:       bytes.NewReader(file),
	})

	if err != nil {
		return err
	}

	info.Parts.Parts = append(info.Parts.Parts, s3Types.CompletedPart{
		ETag:       resp.ETag,
		PartNumber: &partNumber,
	})

	au.mu.Lock()
	au.sessions[sessionInfo.ID] = info
	au.mu.Unlock()

	return nil
}

func (au *AwsUploader) CompletePartUpload(sessionInfo *types.SessionInfo, file []byte) (*UploadAck, error) {
	err := au.HandlePartUpload(sessionInfo, file)

	if err != nil {
		return nil, err
	}

	au.mu.Lock()
	info, isExists := au.sessions[sessionInfo.ID]
	au.mu.Unlock()

	if !isExists {
		return nil, fmt.Errorf("no such session exists")
	}

	resp, err := au.s3Client.CompleteMultipartUpload(context.TODO(), &s3.CompleteMultipartUploadInput{
		Bucket:          info.UploadInfo.Bucket,
		Key:             info.UploadInfo.Key,
		UploadId:        info.UploadInfo.UploadId,
		MultipartUpload: &info.Parts,
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println(resp)

	au.mu.Lock()
	delete(au.sessions, sessionInfo.ID)
	au.mu.Unlock()

	return &UploadAck{
		Provider: types.UploaderAws,
		ProviderMetadata: map[string]interface{}{
			"bucket":   au.config.Store.AWS.Bucket,
			"key":      sessionInfo.FileName,
			"checksum": resp.ChecksumSHA256,
			"etag":     resp.ETag,
		},
		Checksum: "",
		FilePath: "",
	}, nil
}

func (au *AwsUploader) CancelUploadSession(sessionInfo *types.SessionInfo) error {
	au.mu.Lock()
	session, isExists := au.sessions[sessionInfo.ID]
	au.mu.Unlock()

	if !isExists {
		return fmt.Errorf("no such upload session found")
	}

	_, err := au.s3Client.AbortMultipartUpload(context.TODO(), &s3.AbortMultipartUploadInput{
		Bucket:   session.UploadInfo.Bucket,
		Key:      session.UploadInfo.Bucket,
		UploadId: session.UploadInfo.UploadId,
	})

	if err != nil {
		return err
	}

	au.mu.Lock()
	delete(au.sessions, sessionInfo.ID)
	au.mu.Unlock()

	return nil
}

func (au *AwsUploader) DeleteUpload(uploadInfo *db.DeleteDocumentsByIdsRow) error {
	var metadata Metadata

	if err := utils.GetJsonStruct(uploadInfo.ProviderMetadata, &metadata); err != nil {
		return fmt.Errorf("invalid metadata for the provider")
	}

	fmt.Println(metadata.Bucket)
	fmt.Println(metadata.Key)

	_, err := au.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &metadata.Bucket,
		Key:    &metadata.Key,
	})

	if err != nil {
		return err
	}

	return nil
}
