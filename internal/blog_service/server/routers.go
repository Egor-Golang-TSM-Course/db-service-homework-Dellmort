package server

import "github.com/go-chi/chi/v5/middleware"

func (s *Server) routers() {
	s.router.Use(middleware.Logger)

	s.router.Get("/", s.Default())
	s.router.Post("/users/register", s.RegisterUser())
	s.router.Post("/users/login", s.LoginUser())
}
