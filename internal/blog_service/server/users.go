package server

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"homework/internal/models"
	"log/slog"
	"net/http"
	"time"
)

var (
	ErrJsonUnmarshal = errors.New("json Unmarshaler error")
	ErrincorrectData = errors.New("login or password incorrect")
	ErrUndefinedUser = errors.New("user not register")
)

func generateSHA1(password string) string {
	h := sha1.New()
	h.Write([]byte(password))

	return fmt.Sprintf("%x", h.Sum(nil))
}

// TODO: add validions
func (s *Server) registerUser() func(w http.ResponseWriter, r *http.Request) {
	const f = "server.RegisterUser"

	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{
			Created: time.Now(),
		}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusInternalServerError, ErrJsonUnmarshal)
			return
		}

		err = s.validator.Struct(user)
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusBadRequest, err)
			return
		}
		user.Password = generateSHA1(user.Password)

		id, err := s.storage.CreateUser(r.Context(), &user)
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, http.StatusCreated, id)
	}
}

func (s *Server) loginUser() func(w http.ResponseWriter, r *http.Request) {
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
		user, err := s.storage.GetUserByLogin(r.Context(), resp.Login)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				s.error(w, http.StatusBadGateway, ErrUndefinedUser)
				return
			}
			s.error(w, http.StatusBadGateway, err)
			return
		}

		if generateSHA1(resp.Password) != user.Password {
			s.error(w, http.StatusNotAcceptable, ErrincorrectData)
			return
		}

		jwtToken, err := s.TokenManager.GenerateJWTToken(user.ID, user.Login)
		if err != nil {
			s.log.Error(err.Error())
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		// TODO: refresh token
		s.respond(w, http.StatusAccepted, map[string]string{
			"access_token": jwtToken,
		})
	}
}
