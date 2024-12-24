package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/utils"
)

func (s *Services) serviceCreateSession(c echo.Context) error {
	payload := new(types.CreateSessionPayload)

	authInfo := c.Get(string(types.ContextAuthInfo)).(*authInfo)

	if err := utils.GetValidatedPayload(c, payload); err != nil {
		return err
	}

	sessionKey := uuid.New()

	if s.cfg.Redis != "" {
		_, err := s.redisCache.CreateSession(sessionKey.String(), &types.SessionInfo{
			FileSize:  payload.FileSize,
			ChunkSize: payload.ChunkSize,
			MaxChunk:  payload.MaxChunk,
			FileName:  payload.FileName,
			OwnerID:   authInfo.ID,
			TempPath:  sessionKey.String() + "/" + payload.FileName,
		})

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
	}

	sessionToken, err := utils.GetSessionToken(&utils.SessionClaims{
		SessionID: sessionKey.String(),
	}, "secret key")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if err := os.MkdirAll(s.cfg.Store.Local.Temp+"/"+sessionKey.String()+payload.FileName, os.ModeDir); err != nil {
		return err
	}

	s.tasks.AddTask(sessionKey.String(), time.Duration(time.Second*5), func(ctx context.Context) {
		if err := os.RemoveAll(s.cfg.Store.Local.Temp + "/" + sessionKey.String() + payload.FileName); err != nil {
			fmt.Printf(err.Error())
		}
	})

	return c.JSON(http.StatusOK, map[string]string{
		"sessionToken": sessionToken,
	})
}

func (s *Services) registerSessionServicer() {
	group := s.e.Group("/services")

	group.Use(s.middlewareAuthenticate)

	group.POST("", s.serviceCreateSession)
}
