package uploader

import (
	"os"

	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/utils"
)

type LocalUploader struct {
	tempPath  string
	storePath string
	fio       *utils.FileIO
}

type LocalUploaderCompleteRes struct {
	FilePath string
	Checkum  string
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

func (l *LocalUploader) Upload(sessionInfo *types.SessionInfo, file []byte) (*UploadAck, error) {
	path, err := l.fio.WriteToStoreOnly(sessionInfo.FileName, file)

	return &UploadAck{
		Provider:         types.UploaderLocal,
		ProviderMetadata: map[string]string{},
		Checksum:         "",
		FilePath:         path,
	}, err
}

func (l *LocalUploader) HandlePartUpload(sessionInfo *types.SessionInfo, file []byte) error {
	return l.fio.AddChunk(sessionInfo.TempPath, int(sessionInfo.CurrentChunk)+1, file)
}

func (l *LocalUploader) CompletePartUpload(sessionInfo *types.SessionInfo, file []byte) (*UploadAck, error) {
	filePath, checksum, err := l.fio.MergeChunksAndWriteToStore(sessionInfo.TempPath, sessionInfo.FileName, sessionInfo.MaxChunk, file)

	if err != nil {
		return nil, err
	}

	return &UploadAck{
		Provider:         types.UploaderLocal,
		ProviderMetadata: map[string]string{},
		Checksum:         checksum,
		FilePath:         filePath,
	}, nil
}

func (l *LocalUploader) CancelUploadSession(sessionInfo *types.SessionInfo) error {
	return os.RemoveAll(sessionInfo.TempPath)
}

func (l *LocalUploader) DeleteUpload(uploadInfo *db.DeleteDocumentsByIdsRow) error {
	return l.fio.DeleteFile(uploadInfo.FilePath)
}
