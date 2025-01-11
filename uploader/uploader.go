package uploader

import (
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/types"
)

type IUploader interface {
	CreateChunkUploadSession(sessionInfo *types.SessionInfo) error

	Upload(sessionInfo *types.SessionInfo, file []byte) error
	HandlePartUpload(sessionInfo *types.SessionInfo, file []byte) error

	CancelUploadSession(sessionInfo *types.SessionInfo) error
}

type Uploader struct {
	localUploader IUploader
}

func NewUploader(cfg *config.LancerConfig) *Uploader {
	localUploader := NewLocalUploader(cfg.Store.Local.Temp, cfg.Store.Local.Path)

	return &Uploader{
		localUploader: localUploader,
	}
}
