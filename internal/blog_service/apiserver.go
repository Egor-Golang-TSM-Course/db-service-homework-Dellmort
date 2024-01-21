package blogservice

import (
	"homework/internal/blog_service/server"
	"homework/internal/config"
	"homework/internal/logger"
	"homework/storage/postgresql"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Start(cfg *config.Config, validator *validator.Validate) error {
	router := chi.NewRouter()
	database := postgresql.NewPostgreSQLStorage(NewDB(cfg))
	log := logger.Logger(logger.JSON)
	server := server.New(database, router, log, validator)

	return server.Start(cfg)
}

func NewDB(cfg *config.Config) *sqlx.DB {
	const f = "blogservice.NewDB"

	db, err := sqlx.Connect("postgres", cfg.DbConn)
	if err != nil {
		slog.Error(f, slog.String("err", err.Error()))
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		slog.Error(f, slog.String("err", err.Error()))
		os.Exit(1)
	}

	return db
}
