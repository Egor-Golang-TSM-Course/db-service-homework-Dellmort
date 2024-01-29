package server

import (
	"encoding/json"
	"homework/internal/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) addTagsToPost() func(w http.ResponseWriter, r *http.Request) {
	const f = "server.addTagsToPost"

	type response struct {
		Tags string `json:"tags"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.Atoi(chi.URLParam(r, "postId"))
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusBadRequest, ErrInvalidID)
			return
		}
		userID := r.Context().Value(ID("id")).(int)

		post, err := s.storage.GetPostByID(r.Context(), postID)
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		var resp response
		err = json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusBadGateway, ErrJsonUnmarshal)
			return
		}
		if post.OwnerID != userID {
			s.error(w, http.StatusForbidden, ErrNotOwner)
			return
		}

		err = s.storage.CreateTags(r.Context(), &models.Tags{
			Tags:   resp.Tags,
			PostID: postID,
		})
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, http.StatusOK, "ok")
	}
}

func (s *Server) getTagsByPost() func(w http.ResponseWriter, r *http.Request) {
	const f = "server.getTagsByPost"

	return func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.Atoi(chi.URLParam(r, "postId"))
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusBadRequest, ErrInvalidID)
			return
		}

		tags, err := s.storage.GetTagsByPostID(r.Context(), postID)
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, http.StatusOK, tags)
	}
}
