package services

import (
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/types"
)

func (s *Services) serviceGetSettings(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"database":                "postgres",
		"databaseName":            s.cfg.Database.Name,
		"authWebhook":             s.cfg.AuthEndpoint,
		"eventsWebhook":           s.cfg.WebhookEndpoint,
		"webhookSecret":           s.cfg.WebhookSigningSecret,
		"storePath":               s.cfg.Store.Local.Path,
		"tempPath":                s.cfg.Store.Local.Temp,
		"isAwsEnabled":            s.cfg.Store.AWS.Store,
		"awsRegion":               s.cfg.Store.AWS.Region,
		"awsBucket":               s.cfg.Store.AWS.Bucket,
		"isRedis":                 s.cfg.UseRedis,
		"redisServer":             s.cfg.Redis,
		"allowedOrigins":          s.cfg.AllowOrigin,
		"isAuthenticationEnabled": s.cfg.ServerAuth,
		"port":                    s.cfg.Port,
	})
}

func (s *Services) registerSettingsService() {
	group := s.e.Group("/settings")

	group.GET("", s.middlewareAuth([]types.AuthKeys{types.AuthWebToken}, s.serviceGetSettings))
}
