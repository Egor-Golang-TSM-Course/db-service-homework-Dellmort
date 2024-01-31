package blogservice

import (
	"context"
	"homework/internal/blog_service/server"
	"homework/internal/config"
	"homework/internal/logger"
	"homework/storage/postgresql"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

func Start(cfg *config.Config, validator *validator.Validate) error {
	var tokenManager server.TokenManager = server.NewJWTManager(cfg.JwtSecretKey)
	ctx := context.Background()

	router := chi.NewRouter()
	database := postgresql.NewPostgreSQLStorage(NewDB(ctx, cfg))
	log := logger.Logger(logger.JSON)
	serve := server.New(cfg, database, router, log, validator, tokenManager)

	return serve.Start(tokenManager)
}

func NewDB(ctx context.Context, cfg *config.Config) *pgx.Conn {
	const f = "blogservice.NewDB"
	db, err := pgx.Connect(ctx, cfg.DbConn)
	if err != nil {
		slog.Error(f, slog.String("err", err.Error()))
		os.Exit(1)
	}

	err = db.Ping(ctx)
	if err != nil {
		slog.Error(f, slog.String("err", err.Error()))
		os.Exit(1)
	}

	return db
}
