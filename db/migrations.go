package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func openDB(databaseUrl string) *sql.DB {
	db, err := sql.Open("postgres", databaseUrl)

	if err != nil {
		log.Fatalf("[Lancer Error] Failed to open database (%v)", err)
	}

	return db
}

func RunMigration(databaseUrl string) {
	source, err := iofs.New(migrationFiles, "migrations")

	if err != nil {
		log.Fatalf("[Lancer Error] Failed to create source from embeded files (%v)", err)
	}

	dbDriver, err := postgres.WithInstance(openDB(databaseUrl), &postgres.Config{})

	if err != nil {
		log.Fatalf("[Lancer Error] Failed to create database driver for migrations (%v)", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", dbDriver)

	if err != nil {
		log.Fatalf("[Lancer Error] Database migration initialization failed. (%v)", err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("[Lancer Error] Database migration failed: (%v)", err)
	}

	fmt.Println("[Lancer Info] Migration Applied Successfully")

	version, _, _ := m.Version()

	fmt.Printf("Current Version: %d\n", version)
}
