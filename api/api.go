package api

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/weekend-dev-labs/lancer/cache"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/dashboard"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/db/repo"
	"github.com/weekend-dev-labs/lancer/services"
	"github.com/weekend-dev-labs/lancer/utils"
)

func StartServer(cfg *config.LancerConfig, db *db.Queries, cache *cache.Cache, logger *logrus.Logger) {
	e := echo.New()
	e.HideBanner = true

	// _, err := fs.Sub(dashboardFiles, "../dashboard/dist")

	// if err != nil {
	// 	log.Fatal(err)
	// }

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

	e.GET("/*", func(c echo.Context) error {
		path := c.Request().URL.Path

		if !strings.Contains(path, ".") {
			path = "index.html"
		}

		data, err := dashboard.Open(path)

		if err != nil {
			log.Print(err)
			return c.NoContent(http.StatusNotFound)
		}

		mimeType := utils.GetMimetypeByPath(path)

		c.Response().Header().Set(echo.HeaderContentType, mimeType)

		return c.Stream(http.StatusOK, mimeType, bytes.NewReader(data))
	})

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
