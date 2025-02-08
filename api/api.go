package api

import (
	"fmt"
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
)

func StartServer(cfg *config.LancerConfig, db *db.Queries, cache *cache.Cache, logger *logrus.Logger) {
	e := echo.New()
	e.HideBanner = true

	e.Static("/media", cfg.Store.Local.Path)

	e.Validator = &services.LancerValidator{
		Validator: validator.New(),
	}

	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowOrigin,
	}))

	newRepo := repo.NewRepo(db, cache, cfg)

	services.RegisterServices(e.Group("/api"), db, cfg, cache, newRepo, logger)

	startLog := fmt.Sprintf(`                                                                              
   __                        
  / /  ___ ____  _______ ____
 / /__/ _ \/ _ \/ __/ -_) __/
/____/\_,_/_//_/\__/\__/_/    %s

Thanks for using Lancer !!
                             
	`, config.Version)

	fmt.Println(startLog)

	if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("[Lancer Error] Failed to start HTTP Server (%v)", err)
	}
}
