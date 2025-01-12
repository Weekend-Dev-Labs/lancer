package services

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/weekend-dev-labs/lancer/cache"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/db/repo"
	"github.com/weekend-dev-labs/lancer/uploader"
)

type Services struct {
	e           *echo.Group
	db          *db.Queries
	cfg         *config.LancerConfig
	redisCache  *cache.Cache
	tasks       *TaskManager
	repo        *repo.Repo
	webhook     *Webhook
	logger      *logrus.Logger
	uploader    *ServiceUploader
	appUploader *uploader.Uploader
}

type ServiceUploader struct {
	Aws *uploader.AwsUploader
}

func RegisterServices(e *echo.Group, db *db.Queries, cfg *config.LancerConfig, redisCache *cache.Cache, repo *repo.Repo, logger *logrus.Logger) {

	taskManager := NewTaskManager(repo)
	webhook := NewWebhookNotifier(cfg.WebhookEndpoint)

	// appUploader := uploader.

	appUploader := uploader.NewUploader(cfg)

	services := Services{
		e:           e,
		db:          db,
		cfg:         cfg,
		redisCache:  redisCache,
		tasks:       taskManager,
		repo:        repo,
		webhook:     webhook,
		logger:      logger,
		appUploader: appUploader,
	}

	// registering the services
	services.registerSessionServicer()
	services.registerUploaderService()
	services.registerAdminService()
	services.registerMetricsService()
}
