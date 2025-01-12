package uploader

import (
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/types"
)

type IUploader interface {
	CreateChunkUploadSession(sessionInfo *types.SessionInfo) error

	Upload(sessionInfo *types.SessionInfo, file []byte) (*UploadAck, error)
	CompletePartUpload(sessionInfo *types.SessionInfo, file []byte) (*UploadAck, error)

	HandlePartUpload(sessionInfo *types.SessionInfo, file []byte) error

	CancelUploadSession(sessionInfo *types.SessionInfo) error

	DeleteUpload(uploadInfo *db.DeleteDocumentsByIdsRow) error
}

type Uploader struct {
	localUploader IUploader
	awsUploader   IUploader
}

type UploadAck struct {
	Provider         types.UploaderProvider
	ProviderMetadata interface{}
	Checksum         string
	FilePath         string
}

func NewUploader(cfg *config.LancerConfig) *Uploader {
	localUploader := NewLocalUploader(cfg.Store.Local.Temp, cfg.Store.Local.Path)
	awsUploader := NewAwsUploader(cfg)

	return &Uploader{
		localUploader: localUploader,
		awsUploader:   awsUploader,
	}
}

func (u *Uploader) GetUploaderByType(provider types.UploaderProvider) IUploader {
	switch provider {
	case types.UploaderAws:
		return u.awsUploader

	case types.UploaderLocal:
		return u.localUploader
	}

	return nil
}
