package server

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"homework/internal/models"
	"net/http"
	"time"
)

var (
	ErrJsonUnmarshall = errors.New("json unmarshaller error")
	ErrincorrectData  = errors.New("login or password incorrect")
)

func generateSHA1(password string) string {
	h := sha1.New()
	h.Write([]byte(password))

	return fmt.Sprintf("%x", h.Sum(nil))
}

// TODO: add validions
func (s *Server) RegisterUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			s.log.Error(err.Error())
			s.error(w, http.StatusInternalServerError, ErrJsonUnmarshall)
			return
		}

		err = s.validator.Struct(user)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}
		now := time.Now()
		user.Created = &now
		user.Password = generateSHA1(user.Password)

		id, err := s.storage.CreateUser(r.Context(), &user)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, http.StatusOK, id)
	}
}

// TODO: JWT
func (s *Server) LoginUser() func(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var resp response

		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			s.error(w, http.StatusBadGateway, err)
			return
		}
		user, err := s.storage.GetUser(r.Context(), resp.Login)
		if err != nil {
			s.error(w, http.StatusBadGateway, err)
			return
		}
		if generateSHA1(resp.Password) != user.Password {
			s.error(w, http.StatusNotAcceptable, ErrincorrectData)
			return
		}

		// JWT TOKEN GENERATE
		s.respond(w, http.StatusAccepted, user)
	}
}
