package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/types"
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

	if sessionInfo.MaxChunk == 1 {

		if err := s.fio.WriteToStoreOnly(sessionInfo.FileName, fileData); err != nil {
			return err
		}

		s.tasks.CancelTask(session.SessionID)

		return c.JSON(http.StatusAccepted, map[string]string{
			"message": "file uploaded",
		})
	}

	if sessionInfo.CurrentChunk+1 == int64(sessionInfo.MaxChunk) {

		filePath, checksum, err := s.fio.MergeChunksAndWriteToStore(sessionInfo.TempPath, sessionInfo.FileName, sessionInfo.MaxChunk, fileData)

		if err != nil {
			return nil
		}

		if err := s.tasks.CancelWithBaseTask(sessionInfo.TempPath, session.SessionID); err != nil {
			return err
		}

		ext := filepath.Ext(sessionInfo.FileName)
		mimeType := mime.TypeByExtension(ext)

		_, err = s.db.CreateUploadedFile(context.TODO(), db.CreateUploadedFileParams{
			FileName: sessionInfo.FileName,
			FilePath: filePath,
			FileSize: sessionInfo.FileSize,
			FileType: pgtype.Text{
				String: mimeType,
				Valid:  true,
			},
			UploadedBy: sessionInfo.OwnerID,
			Provider:   "local",
			Checksum: pgtype.Text{
				String: checksum,
				Valid:  true,
			},
		})

		if err != nil {
			return err
		}

		return c.JSON(http.StatusAccepted, map[string]string{
			"success": checksum,
		})
	}

	if err := s.fio.AddChunk(sessionInfo.TempPath, int(sessionInfo.CurrentChunk)+1, fileData); err != nil {
		return err
	}

	if err := s.repo.UpdateSessionById(session.SessionID, &types.SessionInfo{
		CurrentChunk: int64(chunkCount),
		FileName:     sessionInfo.FileName,
		TempPath:     sessionInfo.TempPath,
	}); err != nil {
		return err
	}

	s.tasks.ExtendDuration(session.SessionID, time.Duration(time.Minute*2))

	return c.JSON(http.StatusOK, map[string]string{
		"serverChecksum": serverChecksum,
	})
}

func (s *Services) registerUploaderService() {
	group := s.e.Group("/upload")

	group.POST("", s.middlewareSessionAuthenticator(s.serviceHandlerChunkUploader))
}
