package services

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/utils"
)

func (s *Services) serviceHandlerChunkUploader(c echo.Context) error {

	session := c.Get(string(types.ContextSessionInfo)).(*utils.SessionClaims)

	var sessionInfo *types.SessionInfo

	if s.cfg.Redis != "" {
		info, err := s.redisCache.GetSessionInfo(session.SessionID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "no session found",
			})
		}

		sessionInfo = info
	} else {

		pgUuid, err := uuid.Parse(session.SessionID)

		info, err := s.db.FindSessionById(context.TODO(), pgtype.UUID{
			Bytes: pgUuid,
			Valid: true,
		})

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "no session found",
			})
		}

		sessionInfo = &types.SessionInfo{
			FileSize:     info.FileSize,
			ChunkSize:    info.ChunkSize,
			MaxChunk:     info.MaxChunk,
			FileName:     info.FileName.String,
			TempPath:     info.TempPath.String,
			CurrentChunk: 0,
			OwnerID:      info.OwnerID.String,
		}
	}

	return c.JSON(http.StatusOK, sessionInfo)
}

func (s *Services) registerUploaderService() {
	group := s.e.Group("/upload")

	group.POST("", s.middlewareSessionAuthenticator(s.serviceHandlerChunkUploader))
}
