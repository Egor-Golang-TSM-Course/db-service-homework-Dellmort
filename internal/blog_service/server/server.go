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
	storage   storage.Storage
	router    chi.Router
	log       *slog.Logger
	validator *validator.Validate
}

func New(storage storage.Storage, router chi.Router, log *slog.Logger, validator *validator.Validate) *Server {
	return &Server{
		storage:   storage,
		router:    router,
		log:       log,
		validator: validator,
	}
}

func (s *Server) Start(cfg *config.Config) error {
	s.log.Info("server started...")
	s.routers()

	return http.ListenAndServe(":8080", s.router)
}
