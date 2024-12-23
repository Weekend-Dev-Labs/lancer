package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/utils"
)

func (s *Services) serviceCreateSession(c echo.Context) error {
	payload := new(types.CreateSessionPayload)

	if err := utils.GetValidatedPayload(c, payload); err != nil {
		return err
	}

	if s.cfg.Redis != "" {
		// session := s.
	}

	return c.JSON(http.StatusOK, payload)
}

func (s *Services) registerSessionServicer() {
	group := s.e.Group("/services")

	group.Use(s.middlewareAuthenticate)

	group.POST("", s.serviceCreateSession)
}
