package api

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/weekend-dev-labs/lancer/cache"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/db/repo"
	"github.com/weekend-dev-labs/lancer/services"
	"github.com/weekend-dev-labs/lancer/uploader"
)

func StartServer(cfg *config.LancerConfig, db *db.Queries, cache *cache.Cache, logger *logrus.Logger) {
	e := echo.New()

	e.Validator = &services.LancerValidator{
		Validator: validator.New(),
	}

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	newRepo := repo.NewRepo(db, cache, cfg)

	awsUploader := uploader.NewAwsUploader(cfg)

	services.RegisterServices(e.Group("/api"), db, cfg, cache, newRepo, logger, &services.ServiceUploader{
		Aws: awsUploader,
	})

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("[Lancer Error] Failed to start HTTP Server (%v)", err)
	}
}
