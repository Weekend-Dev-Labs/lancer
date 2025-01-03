package services

import (
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/types"
)

func (s *Services) serviceGetMetrics(c echo.Context) error {
	return nil
}

func (s *Services) registerMetricsService() {
	group := s.e.Group("/metrics")

	group.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return s.middlewareAuth([]types.AuthKeys{types.AuthWebToken}, next)
	})

	group.GET("", s.serviceGetMetrics)
}
