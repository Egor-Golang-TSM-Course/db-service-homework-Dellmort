package server

import (
	"encoding/json"
	"errors"
	"homework/internal/models"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

var (
	ErrInvalidID = errors.New("id not number")
	ErrNotFound  = errors.New("not found")
)

func (s *Server) createPost() func(w http.ResponseWriter, r *http.Request) {
	const f = "server.CreatePost"

	return func(w http.ResponseWriter, r *http.Request) {
		post := models.Post{
			Created: time.Now(),
		}

		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			slog.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusInternalServerError, err)
			return
		}
		err = s.validator.Struct(post)
		if err != nil {
			slog.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusBadRequest, err)
			return
		}

		post.OwnerID = r.Context().Value(ID("id")).(int)
		id, err := s.storage.CreatePost(r.Context(), &post)
		if err != nil {
			slog.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, http.StatusOK, id)
	}
}

// TODO: Добавить фильтрацию по тегам или дате
func (s *Server) getAllPosts() func(w http.ResponseWriter, r *http.Request) {
	const f = "server.GetAllPosts"
	type response struct {
		Tags []string   `json:"tags" validate:"required"`
		Date *time.Time `json:"date" validate:"required"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var resp response
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusBadRequest, ErrJsonUnmarshal)
			return
		}

		err := s.validator.Struct(resp)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		posts, err := s.storage.GetPosts(r.Context(), resp.Tags, resp.Date)
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		if len(posts) == 0 {
			s.error(w, http.StatusNotFound, ErrNotFound)
			return
		}

		postsMap := make(map[int]*models.Post, 0)
		for i, post := range posts {
			postsMap[i+1] = post
		}

		s.respond(w, http.StatusOK, postsMap)
	}
}

func (s *Server) getPostsById() func(w http.ResponseWriter, r *http.Request) {
	const f = "server.GetAllPosts"
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		postID, err := strconv.Atoi(id)
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusBadRequest, ErrInvalidID)
			return
		}

		// TODO: validation

		post, err := s.storage.GetPostByID(r.Context(), postID)
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusBadRequest, err)
			return
		}

		s.respond(w, http.StatusOK, *post)
	}
}

func (s *Server) updatePost() func(w http.ResponseWriter, r *http.Request) {
	const f = "server.updatePost"
	type response struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusBadRequest, ErrInvalidID)
			return
		}
		resp := response{}
		err = json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			s.error(w, http.StatusBadRequest, ErrJsonUnmarshal)
			return
		}
		post := models.NewPost()
		post.Name = resp.Name
		post.Message = resp.Message
		post.ID = postID

		err = s.storage.UpdatePost(r.Context(), post)
		if err != nil {
			s.error(w, http.StatusInternalServerError, nil)
			return
		}

		s.respond(w, http.StatusOK, "updated")
	}
}

func (s *Server) deletePost() func(w http.ResponseWriter, r *http.Request) {
	const f = "server.deletePost"
	return func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusBadRequest, ErrInvalidID)
			return
		}

		err = s.storage.DeletePost(r.Context(), postID)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, http.StatusOK, "deleted")
	}
}
