package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/utils"
)

func (s *Services) serviceCreateSession(c echo.Context) error {
	payload := new(types.CreateSessionPayload)

	authInfo := c.Get(string(types.ContextAuthInfo)).(*authInfo)

	if err := utils.GetValidatedPayload(c, payload); err != nil {
		return err
	}

	uploadHandler := s.appUploader.GetUploaderByType(payload.Provider)

	if uploadHandler == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid provider to create session",
		})
	}

	// var isAwsUploader bool

	// switch payload.Provider {
	// case types.UploaderAws:
	// 	isAwsUploader = true
	// 	if !s.cfg.Store.AWS.Store {
	// 		return c.JSON(http.StatusBadRequest, map[string]string{
	// 			"error": "aws is not configured to handle uploads",
	// 		})
	// 	}
	// }

	session, err := s.repo.CreateSession(&types.SessionInfo{
		FileSize:     payload.FileSize,
		ChunkSize:    payload.ChunkSize,
		MaxChunk:     payload.MaxChunk,
		FileName:     payload.FileName,
		OwnerID:      authInfo.ID,
		CurrentChunk: 0,
		Provider:     payload.Provider,
	})

	fmt.Printf(string(session.Provider))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	sessionToken, err := utils.GetSessionToken(&utils.SessionClaims{
		SessionID: session.ID,
	}, "secret key")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if err := uploadHandler.CreateChunkUploadSession(session); err != nil {
		return err
	}

	s.tasks.AddTask(session.ID, time.Duration(time.Second*300), func(ctx context.Context) {
		uploadHandler.CancelUploadSession(session)
	})
	// 	if isAwsUploader {
	// 		s.uploader.Aws.AbortMultipartUpload(session.ID)
	// 	} else {
	// 		if err := os.RemoveAll(dirPath); err != nil {
	// 			fmt.Printf(err.Error())
	// 		}
	// 	}
	// })

	go s.webhook.SendEvent(EventSessionCreate, session)

	return c.JSON(http.StatusOK, map[string]string{
		"sessionToken": sessionToken,
	})
}

func (s *Services) serviceEndSession(c echo.Context) error {
	payload := new(types.SessionTokenPayload)

	if err := utils.GetValidatedPayload(c, payload); err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error": "session_token required",
		})
	}

	sessionInfo, err := utils.GetSessionInfo(payload.SessionToken, "secret key")

	if err != nil {
		return err
	}

	var session *types.SessionInfo

	if s.cfg.Redis != "" {
		session, err = s.redisCache.GetSessionInfo(sessionInfo.SessionID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid session or session not found",
			})
		}
	} else {
		sessionUuid, err := uuid.Parse(sessionInfo.SessionID)

		if err != nil {
			return err
		}

		dbSession, err := s.db.FindSessionById(context.Background(), pgtype.UUID{
			Bytes: sessionUuid,
			Valid: true,
		})

		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid session or session not found",
			})
		}

		session = &types.SessionInfo{
			FileSize:     dbSession.FileSize,
			ChunkSize:    dbSession.ChunkSize,
			MaxChunk:     dbSession.MaxChunk,
			FileName:     dbSession.FileName.String,
			TempPath:     dbSession.TempPath.String,
			OwnerID:      dbSession.OwnerID.String,
			CurrentChunk: 0,
		}
	}

	if err := BaseTask(session.TempPath, context.Background()); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	s.tasks.Execute(sessionInfo.SessionID)

	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "session ended",
	})

}

func (s *Services) serviceGetSessions(c echo.Context) error {

	sessions := s.repo.GetSessions(&db.PaginateSessionsParams{
		Limit:  int32(10),
		Offset: int32(0),
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sessions": sessions,
	})
}

func (s *Services) registerSessionServicer() {
	group := s.e.Group("/sessions")

	// group.Use(s.middlewareAuthenticate)

	group.GET("", s.middlewareAuth([]types.AuthKeys{types.AuthWebToken}, s.serviceGetSessions))

	group.POST("", s.middlewareAuth([]types.AuthKeys{types.AuthClientServerToken}, s.serviceCreateSession))
	group.POST("/end", s.middlewareAuth([]types.AuthKeys{types.AuthClientServerToken}, s.serviceEndSession))
}
