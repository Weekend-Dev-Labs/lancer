package services

import (
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
)

type Services struct {
	e   *echo.Group
	db  *db.Queries
	cfg *config.LancerConfig
}

func RegisterServices(e *echo.Group, db *db.Queries, cfg *config.LancerConfig) {
	services := Services{
		e:   e,
		db:  db,
		cfg: cfg,
	}

	// registering the services
	services.registerSessionServicer()
}
