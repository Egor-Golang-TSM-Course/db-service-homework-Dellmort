package server

import (
	"encoding/json"
	"errors"
	"homework/internal/models"
	"homework/internal/pkg"
	"log/slog"
	"net/http"
	"time"
)

var (
	ErrJsonUnmarshal = errors.New("json Unmarshaler error")
	ErrIncorrectData = errors.New("login or password incorrect")
	ErrUndefinedUser = errors.New("user not register")
)

func (s *Server) registerUser() func(w http.ResponseWriter, r *http.Request) {
	const f = "server.RegisterUser"
	type response struct {
		Name     string `json:"name" validate:"required,gte=4"`
		Login    string `json:"login" validate:"required,gte=5"`
		Password string `json:"password" validate:"required,gte=6"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var resp response
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			s.log.Error(f, slog.String("ERR", err.Error()))
			s.error(w, http.StatusInternalServerError, ErrJsonUnmarshal)
			return
		}

		// validate to struct
		err = s.validator.Struct(resp)
		if err != nil {
			s.log.Debug(f, slog.String("DEBUG", err.Error()))
			s.error(w, http.StatusBadRequest, err)
			return
		}

		user := models.User{
			Name:     resp.Name,
			Login:    resp.Login,
			Password: pkg.GenetareSHA1(resp.Password),
			Created:  time.Now(),
		}
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

		if pkg.GenetareSHA1(resp.Password) != user.Password {
			s.error(w, http.StatusNotAcceptable, ErrIncorrectData)
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
