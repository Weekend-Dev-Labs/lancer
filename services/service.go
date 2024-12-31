package services

import (
	"github.com/labstack/echo/v4"
	"github.com/weekend-dev-labs/lancer/cache"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/db/repo"
	"github.com/weekend-dev-labs/lancer/utils"
)

type Services struct {
	e          *echo.Group
	db         *db.Queries
	cfg        *config.LancerConfig
	redisCache *cache.Cache
	tasks      *TaskManager
	repo       *repo.Repo
	fio        *utils.FileIO
	webhook    *Webhook
}

func RegisterServices(e *echo.Group, db *db.Queries, cfg *config.LancerConfig, redisCache *cache.Cache, repo *repo.Repo) {

	taskManager := NewTaskManager(repo)
	webhook := NewWebhookNotifier(cfg.WebhookEndpoint)
	fio := utils.NewFileIO(cfg.Store.Local.Path, cfg.Store.Local.Temp)

	services := Services{
		e:          e,
		db:         db,
		cfg:        cfg,
		redisCache: redisCache,
		tasks:      taskManager,
		repo:       repo,
		fio:        fio,
		webhook:    webhook,
	}

	// registering the services
	services.registerSessionServicer()
	services.registerUploaderService()
	services.registerAdminService()
}
