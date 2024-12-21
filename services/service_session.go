package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/types"
)

func (s *Services) serviceCreateSession(c echo.Context) error {

	// _, _ := c.Get(string(types.ContextAuthInfo)).(*authInfo)

	payload := new(types.CreateSessionPayload)

	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "invalid payload")
	}

	if err := c.Validate(payload); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, payload)
}

func (s *Services) registerSessionServicer() {
	group := s.e.Group("/services")

	group.Use(s.middlewareAuthenticate)

	group.POST("", s.serviceCreateSession)
}
