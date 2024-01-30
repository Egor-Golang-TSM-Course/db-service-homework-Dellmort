package main

import (
	"errors"
	"flag"
	"fmt"
	blogservice "homework/internal/blog_service"
	"homework/internal/config"
	"homework/internal/pkg/migrations"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
)

const (
	migrationFileURL = "file://migrations"
)

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDYxNTA2MzQsImlkIjozfQ.eu-7Fcbgnedj6Hp2bUTmR47uNd269DsuyBj7t1578vI

func main() {
	configPath := flag.String("cfg", "config/config.yml", "path to file config")
	flag.Parse()
	if *configPath == "" {
		panic("nil config file")
	}

	validator := validator.New()
	cfg := config.MustConfig(*configPath, validator)

	if err := migrations.Migration(migrationFileURL, cfg.DbConn); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error(fmt.Sprintf("UP migrate to database ERROR %v", err))
	}

	blogservice.Start(cfg, validator)
}
