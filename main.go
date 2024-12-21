package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"github.com/weekend-dev-labs/lancer/api"
	"github.com/weekend-dev-labs/lancer/config"
	"github.com/weekend-dev-labs/lancer/db"
)

var logger = logrus.New()

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg := config.ParseFlags()

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

	api.StartServer(cfg, query)
}
