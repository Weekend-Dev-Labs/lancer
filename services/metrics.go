package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/types"
)

func (s *Services) serviceGetMetrics(c echo.Context) error {
	metrics, _ := s.db.GetFirstCreatedMetrics(context.TODO())

	metricsDecorder := json.NewDecoder(bytes.NewReader(metrics.FilesByMimetype))

	var mimetypeMetrics map[string]interface{}

	_ = metricsDecorder.Decode(&mimetypeMetrics)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"metrics": map[string]interface{}{
			"ID":                metrics.ID,
			"TotalFileSize":     metrics.TotalFileSize,
			"TotalFileCount":    metrics.TotalFileCount,
			"FilesByMimetype":   mimetypeMetrics,
			"LargestFileSize":   metrics.LargestFileSize,
			"SmallestFileSize":  metrics.AverageFileSize,
			"TotalDeletedFiles": metrics.TotalDeletedFiles,
			"LastUpdated":       metrics.LastUpdated,
		},
	})
}

func (s *Services) registerMetricsService() {
	group := s.e.Group("/metrics")

	group.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return s.middlewareAuth([]types.AuthKeys{types.AuthWebToken}, next)
	})

	group.GET("", s.serviceGetMetrics)
}
