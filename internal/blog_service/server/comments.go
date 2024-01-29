package server

import (
	"encoding/json"
	"errors"
	"homework/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *Server) addCommentToPost() func(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Msg string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.Atoi(chi.URLParam(r, "postId"))
		if err != nil {
			s.error(w, http.StatusBadRequest, ErrInvalidID)
			return
		}
		ownerID := r.Context().Value(ID("id")).(int)

		var resp response
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			s.error(w, http.StatusBadRequest, ErrJsonUnmarshal)
			return
		}

		if resp.Msg == "" {
			s.error(w, http.StatusBadRequest, errors.New("message is null"))
			return
		}

		commentID, err := s.storage.CreateComment(r.Context(), &models.Comment{
			Message: resp.Msg,
			OwnerID: ownerID,
			PostID:  postID,
			Created: time.Now(),
		})
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, http.StatusOK, commentID)
	}
}

func (s *Server) getAllCommentsToPost() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.Atoi(chi.URLParam(r, "postId"))
		if err != nil {
			s.error(w, http.StatusBadRequest, ErrInvalidID)
			return
		}

		if postID == 0 {
			s.error(w, http.StatusBadRequest, errors.New("post_id less zero"))
			return
		}

		comments, err := s.storage.GetCommentsByPostID(r.Context(), postID)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		if len(comments) == 0 {
			s.respond(w, http.StatusForbidden, ErrNotFound)
			return
		}

		s.respond(w, http.StatusOK, comments)
	}
}
