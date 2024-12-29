package utils

import (
	"fmt"
	"io"
	"os"
	"time"
)

type FileIO struct {
	local string
	temp  string
}

func NewFileIO(local string, temp string) *FileIO {
	return &FileIO{
		local: local,
		temp:  temp,
	}
}

func (fio *FileIO) AddChunk(path string, chunkCount int, fileData []byte) error {
	chunkPath := fmt.Sprintf("%s/chunk_%d", path, chunkCount)

	return os.WriteFile(chunkPath, fileData, os.ModeAppend)
}

func (fio *FileIO) WriteToStoreOnly(fileName string, data []byte) error {
	filePath := fmt.Sprintf("%s/%d_%s", fio.local, time.Now().Unix(), fileName)

	return os.WriteFile(filePath, data, os.ModeAppend)
}

func (fio *FileIO) MergeChunksAndWriteToStore(path string, fileName string, totalChunks int64, data []byte) (string, error) {
	filePath := fmt.Sprintf("%s/%d_%s", fio.local, time.Now().Unix(), fileName)

	outFile, err := os.Create(filePath)

	if err != nil {
		return "", err
	}

	defer outFile.Close()

	for i := 1; i < int(totalChunks); i++ {
		inFile, err := os.Open(path + fmt.Sprintf("/chunk_%d", i))

		if err != nil {
			return "", err
		}

		_, err = io.Copy(outFile, inFile)

		if err != nil {
			return "", err
		}

		inFile.Close()
	}

	return filePath, nil
}
