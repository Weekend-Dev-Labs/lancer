package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/weekend-dev-labs/lancer/api"
	"github.com/weekend-dev-labs/lancer/cache"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
	"github.com/weekend-dev-labs/lancer/db/repo"
)

var logger = logrus.New()

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg := config.ParseFlags()

	fmt.Println(cfg.Database.Migrate)

	if cfg.Database.Migrate {
		db.RunMigration(cfg.GetDatabaseConnectionString())
		return
	}

	if err := os.MkdirAll(cfg.Store.Local.Path, os.ModeAppend); err != nil {
		log.Fatalf("[Lancer Error] Failed to create store directory : %v", err.Error())
	}

	conn, err := pgxpool.New(context.Background(), cfg.GetDatabaseConnectionString())

	if err != nil {
		log.Fatalf("[Lancer Error] Failed to connect to database (%v)", err)
	}

	defer conn.Close()

	query := db.New(conn)

	redisCache := cache.NewCache(cfg.Redis)

	repo.CreateInitialUser(cfg, query)

	if redisCache == nil {
		log.Printf("[LANCER WARNING] Redis Server Configuration Missing. Using Database to store sessions.")
	}

	api.StartServer(cfg, query, redisCache)
}
