package server

import (
	"homework/internal/config"
	"homework/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	cfg          *config.Config
	storage      storage.Storage
	router       chi.Router
	log          *slog.Logger
	validator    *validator.Validate
	TokenManager TokenManager
}

func New(cfg *config.Config,
	storage storage.Storage,
	router chi.Router,
	log *slog.Logger,
	validator *validator.Validate,
	tokenManager TokenManager) *Server {
	return &Server{
		cfg:          cfg,
		storage:      storage,
		router:       router,
		log:          log,
		validator:    validator,
		TokenManager: tokenManager,
	}
}

func (s *Server) Start(tm TokenManager) error {
	s.log.Info("server started...")
	s.routers(tm)

	return http.ListenAndServe(":8080", s.router)
}
