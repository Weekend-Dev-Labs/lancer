package uploader

import (
	"os"

	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/utils"
)

type LocalUploader struct {
	tempPath  string
	storePath string
	fio       *utils.FileIO
}

func NewLocalUploader(tempPath string, storePath string) *LocalUploader {
	fio := utils.NewFileIO(storePath, tempPath)

	return &LocalUploader{
		tempPath:  tempPath,
		storePath: storePath,
		fio:       fio,
	}
}

func (l *LocalUploader) CreateChunkUploadSession(sessionInfo *types.SessionInfo) error {
	if sessionInfo.MaxChunk > 1 {
		dirPath := sessionInfo.TempPath

		if err := os.MkdirAll(dirPath, os.ModeDir); err != nil {
			return err
		}
	}

	return nil
}

func (l *LocalUploader) Upload(sessionInfo *types.SessionInfo, file []byte) (interface{}, error) {
	err := l.fio.WriteToStoreOnly(sessionInfo.FileName, file)

	return nil, err
}

func (l *LocalUploader) HandlePartUpload(sessionInfo *types.SessionInfo, file []byte) error {
	return l.fio.AddChunk(sessionInfo.TempPath, int(sessionInfo.CurrentChunk)+1, file)
}

// func (l *LocalUploader) CompletePartUpload(sess)

func (l *LocalUploader) CancelUploadSession(sessionInfo *types.SessionInfo) error {
	return os.RemoveAll(sessionInfo.TempPath)
}
