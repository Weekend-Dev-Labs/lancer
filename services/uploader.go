package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/uploader"
	"github.com/weekend-dev-labs/lancer/utils"
)

func (s *Services) serviceHandlerChunkUploader(c echo.Context) error {

	session := c.Get(string(types.ContextSessionInfo)).(*utils.SessionClaims)

	sessionInfo, err := s.repo.GetSessionById(session.SessionID)

	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"error": err.Error(),
		})
	}

	uploadHandler := s.appUploader.GetUploaderByType(sessionInfo.Provider)

	if uploadHandler == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid provider to upload file",
		})
	}

	checksum := c.FormValue("checksum")
	chunk := c.FormValue("chunk")

	if checksum == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid checksum",
		})
	}

	chunkCount, err := strconv.Atoi(chunk)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid chunk count",
		})
	}

	hasher := sha256.New()

	file, err := c.FormFile("file")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	content, err := file.Open()

	if err != nil {
		return err
	}

	fileData, err := io.ReadAll(content)

	if err != nil {
		return err
	}

	hasher.Write(fileData)

	serverChecksum := hex.EncodeToString(hasher.Sum(nil))

	if serverChecksum != checksum {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "checksum mismatched",
		})
	}

	ext := filepath.Ext(sessionInfo.FileName)
	mimeType := mime.TypeByExtension(ext)

	if sessionInfo.MaxChunk == 1 {

		// ack, err := s.fio.WriteToStoreOnly(sessionInfo.FileName, fileData)
		// if err != nil {
		// 	return err
		// }

		// s.tasks.CancelTask(session.SessionID)

		uploadRes, err := uploadHandler.Upload(sessionInfo, fileData)

		if err != nil {
			return err
		}

		b, err := json.Marshal(uploadRes.ProviderMetadata)

		if err != nil {
			fmt.Println("Error during marshalling:", err)
			return err
		}

		_, err = s.db.CreateUploadedFile(context.TODO(), db.CreateUploadedFileParams{
			FileName: sessionInfo.FileName,
			FilePath: uploadRes.FilePath,
			FileSize: sessionInfo.FileSize,
			FileType: pgtype.Text{
				String: mimeType,
				Valid:  true,
			},
			UploadedBy: sessionInfo.OwnerID,
			Provider:   string(sessionInfo.Provider),
			Checksum: pgtype.Text{
				String: checksum,
				Valid:  true,
			},
			ProviderMetadata: b,
		})

		return c.JSON(http.StatusAccepted, map[string]string{
			"message": "file uploaded",
		})
	}

	// if sessionInfo.Provider == types.UploaderAws {
	// 	err := s.uploader.Aws.HandleMultipartUploads(sessionInfo.ID, int32(chunkCount), sessionInfo, bytes.NewReader(fileData))

	// 	if err != nil {
	// 		return err
	// 	}

	// 	if err := s.repo.UpdateSessionById(session.SessionID, &types.SessionInfo{
	// 		CurrentChunk: int64(chunkCount),
	// 		FileName:     sessionInfo.FileName,
	// 		TempPath:     sessionInfo.TempPath,
	// 		Provider:     sessionInfo.Provider,
	// 	}); err != nil {
	// 		return err
	// 	}

	// 	s.tasks.ExtendDuration(session.SessionID, time.Duration(time.Minute*5))

	// 	return c.JSON(http.StatusAccepted, map[string]string{
	// 		"message": "uploaded",
	// 	})
	// }

	if sessionInfo.CurrentChunk+1 == int64(sessionInfo.MaxChunk) {

		uploadRes, err := uploadHandler.CompletePartUpload(sessionInfo, fileData)

		if err != nil {
			return err
		}

		b, err := json.Marshal(uploadRes.ProviderMetadata)

		if err != nil {
			fmt.Println("Error during marshalling:", err)
			return err
		}

		// filePath, checksum, err := s.fio.MergeChunksAndWriteToStore(sessionInfo.TempPath, sessionInfo.FileName, sessionInfo.MaxChunk, fileData)

		// if err != nil {
		// 	return nil
		// }

		// if err := s.tasks.CancelWithBaseTask(sessionInfo.TempPath, session.SessionID); err != nil {
		// 	fmt.Println(err.Error())
		// 	return err
		// }

		ext := filepath.Ext(sessionInfo.FileName)
		mimeType := mime.TypeByExtension(ext)

		file, err := s.db.CreateUploadedFile(context.TODO(), db.CreateUploadedFileParams{
			FileName: sessionInfo.FileName,
			FilePath: uploadRes.FilePath,
			FileSize: sessionInfo.FileSize,
			FileType: pgtype.Text{
				String: mimeType,
				Valid:  true,
			},
			UploadedBy: sessionInfo.OwnerID,
			Provider:   string(sessionInfo.Provider),
			Checksum: pgtype.Text{
				String: uploadRes.Checksum,
				Valid:  true,
			},
			ProviderMetadata: b,
		})

		if err != nil {
			log.Fatalf(err.Error())
			return err
		}

		if err := s.db.IncrementFileCountAndSizeAndMimetype(context.TODO(), db.IncrementFileCountAndSizeAndMimetypeParams{
			TotalFileCount: 1,
			TotalFileSize:  file.FileSize,
			ID:             s.cfg.MetricsID,
			Column4:        []byte(file.FileType.String),
		}); err != nil {

			fmt.Printf("\n\n[LANCER ERROR ] %v\n\n", err.Error())

			s.logger.Log(logrus.ErrorLevel, map[string]string{
				"error": err.Error(),
			})
		}

		s.webhook.SendEvent(EventFileUpload, file)

		return c.JSON(http.StatusAccepted, map[string]string{
			"success": checksum,
		})
	}

	if err := uploadHandler.HandlePartUpload(sessionInfo, fileData); err != nil {
		return err
	}

	if err := s.repo.UpdateSessionById(session.SessionID, &types.SessionInfo{
		CurrentChunk: int64(chunkCount),
		FileName:     sessionInfo.FileName,
		TempPath:     sessionInfo.TempPath,
		Provider:     sessionInfo.Provider,
	}); err != nil {
		return err
	}

	s.tasks.ExtendDuration(session.SessionID, time.Duration(time.Minute*2))

	return c.JSON(http.StatusOK, map[string]string{
		"serverChecksum": serverChecksum,
	})
}

func (s *Services) serviceGetUploads(c echo.Context) error {

	var queryInfo types.UploadQueryInfo

	queryInfo.Limit = 20
	queryInfo.Page = 0

	err := c.Bind(&queryInfo)
	if err != nil {
		return err
	}

	uploads, err := s.db.PaginateUploadedFiles(context.Background(), db.PaginateUploadedFilesParams{
		Limit: int32(queryInfo.Limit), Offset: int32(queryInfo.Page),
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"uploads": uploads,
	})
}

func (s *Services) serviceDeleteUploads(c echo.Context) error {
	payload := new(types.UploadDeletePayload)

	if err := utils.GetValidatedPayload(c, payload); err != nil {
		fmt.Printf(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid payload",
		})
	}

	uploadInfo, err := s.db.DeleteDocumentsByIds(context.Background(), payload.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	var wg sync.WaitGroup
	// var deleteErrors []string

	for _, info := range uploadInfo {
		wg.Add(1)
		go func(info db.DeleteDocumentsByIdsRow) {
			defer wg.Done()

			if err := os.RemoveAll(info.FilePath); err == nil {
				err := s.db.DecrementFileCountAndSizeAndMimetype(context.TODO(), db.DecrementFileCountAndSizeAndMimetypeParams{
					TotalFileCount:  1,
					TotalFileSize:   info.FileSize,
					ID:              s.cfg.MetricsID,
					FilesByMimetype: []byte(info.FileType.String),
				})

				if err != nil {
					s.logger.Log(logrus.ErrorLevel, map[string]string{
						"error": err.Error(),
					})
				}
			}
		}(info)
	}

	wg.Wait()

	s.webhook.SendEvent(EventFileDelete, uploadInfo)

	return c.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("%d files deleted", len(uploadInfo)),
	})
}

func (s *Services) serviceUploadFileTestAws(c echo.Context) error {
	file, err := c.FormFile("file")

	if err != nil {
		return err
	}

	content, err := file.Open()

	if err != nil {
		return err
	}

	res, err := s.uploader.Aws.UploadFullFile(&uploader.UploadFullFileParam{
		Bucket: s.cfg.Store.AWS.Bucket,
		Key:    "lancer-test",
		File:   content,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"info": res,
	})
}

func (s *Services) registerUploaderService() {
	group := s.e.Group("/upload")

	group.POST("/aws", s.serviceUploadFileTestAws)
	group.GET("", s.middlewareAuth([]types.AuthKeys{types.AuthWebToken, types.AuthServerCredentials}, s.serviceGetUploads))
	group.POST("/delete", s.middlewareAuth([]types.AuthKeys{types.AuthWebToken, types.AuthServerCredentials}, s.serviceDeleteUploads))
	group.POST("", s.middlewareAuth([]types.AuthKeys{types.AuthSessionToken}, s.serviceHandlerChunkUploader))
}
