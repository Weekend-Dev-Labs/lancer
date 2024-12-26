package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"github.com/weekend-dev-labs/lancer/api"
	"github.com/weekend-dev-labs/lancer/cache"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
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

	conn, err := pgx.Connect(context.Background(), cfg.GetDatabaseConnectionString())

	if err != nil {
		log.Fatalf("[Lancer Error] Failed to connect to database (%v)", err)
	}

	defer conn.Close(context.Background())

	query := db.New(conn)

	redisCache := cache.NewCache(cfg.Redis)

	if redisCache == nil {
		log.Printf("[LANCER WARNING] Redis Server Configuration Missing. Using Database to store sessions.")
	}

	api.StartServer(cfg, query, redisCache)
}
