package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
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

	return c.JSON(http.StatusOK, sessionInfo)
}

func (s *Services) registerUploaderService() {
	group := s.e.Group("/upload")

	group.POST("", s.middlewareSessionAuthenticator(s.serviceHandlerChunkUploader))
}
