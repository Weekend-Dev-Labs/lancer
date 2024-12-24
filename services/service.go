package services

import (
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/cache"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
)

type Services struct {
	e          *echo.Group
	db         *db.Queries
	cfg        *config.LancerConfig
	redisCache *cache.Cache
	tasks      *TaskManager
}

func RegisterServices(e *echo.Group, db *db.Queries, cfg *config.LancerConfig, redisCache *cache.Cache) {

	taskManager := NewTaskManager()

	services := Services{
		e:          e,
		db:         db,
		cfg:        cfg, 
		redisCache: redisCache,
		tasks:      taskManager,
	}

	// registering the services
	services.registerSessionServicer()
}
