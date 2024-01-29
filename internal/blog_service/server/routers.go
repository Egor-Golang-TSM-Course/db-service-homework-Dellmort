package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	version = "/v1"
	API     = "/api" + version
)

func (s *Server) routers(tm TokenManager) {
	s.router.Use(middleware.Logger)
	s.router.Get("/", s.Default())

	s.router.Route(API, func(r chi.Router) {
		r.Post("/register", s.registerUser())
		r.Post("/login", s.loginUser())

		r.With(s.checkAuthMiddleware).Post("/posts", s.createPost())
		r.With(s.checkAuthMiddleware).
			With(s.checkOwnerPostMiddleware).
			Put("/posts/{id}", s.updatePost())

		r.With(s.checkAuthMiddleware).
			With(s.checkOwnerPostMiddleware).
			Delete("/posts/{id}", s.deletePost())

		r.With(s.checkAuthMiddleware).Post("/posts/{postId}/comments", s.addCommentToPost())

		r.With(s.checkAuthMiddleware).
			Post("/posts/{postId}/tags", s.addTagsToPost())

		r.Get("/posts", s.getAllPosts())
		r.Get("/posts/{id}", s.getPostsById())
		r.Get("/posts/{postId}/comments", s.getAllCommentsToPost())
		r.Get("/posts/{postId}/tags", s.getTagsByPost())
	})
}
