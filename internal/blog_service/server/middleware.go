package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

const authorizationHeader = "Authorization"

type ID string

var (
	ErrAuth         = errors.New("empty auth header")
	ErrInvalidToken = errors.New("invalid token")
	ErrNotOwner     = errors.New("user not owner this post")
)

func (s *Server) checkAuthMiddleware(next http.Handler) http.Handler {
	const f = "server.CheckAuthMiddleware"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(authorizationHeader)
		tokenParts := strings.Fields(token)
		if len(tokenParts) != 2 {
			s.error(w, http.StatusUnauthorized, ErrAuth)
			return
		}

		token = tokenParts[1]
		if token == "" {
			s.log.Debug(f, slog.String("ERR", ErrAuth.Error()))
			s.error(w, http.StatusUnauthorized, ErrAuth)
			return
		}

		claims, err := s.TokenManager.ParseJWTtoken(token)
		if err != nil {
			s.log.Debug(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusUnauthorized, ErrAuth)
			return
		}

		if id, ok := claims["id"]; ok {
			userID := int(id.(float64))
			if userID <= 0 {
				s.log.Debug(f, slog.String("ERR", ErrAuth.Error()))
				s.error(w, http.StatusUnauthorized, ErrInvalidToken)
				return
			}

			ctx := context.WithValue(r.Context(), ID("id"), userID)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) checkOwnerPostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pid := chi.URLParam(r, "id")
		uid := r.Context().Value(ID("id")).(int)

		postID, err := strconv.Atoi(pid)
		if err != nil {
			s.error(w, http.StatusBadRequest, ErrInvalidID)
			return
		}
		post, err := s.storage.GetPostByID(r.Context(), postID)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		if post.OwnerID != uid {
			s.error(w, http.StatusForbidden, ErrNotOwner)
			return
		}

		next.ServeHTTP(w, r)
	})
}
