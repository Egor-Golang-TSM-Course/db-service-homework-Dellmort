package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

func (s *Server) Default() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := os.ReadFile("api_plan.md")
		if err != nil {
			return
		}

		s.error(w, http.StatusOK, errors.New(string(body)))
	}
}

func (s *Server) respond(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data = map[string]any{
		"response": data,
	}

	return json.NewEncoder(w).Encode(data)
}

func (s *Server) error(w http.ResponseWriter, code int, err error) error {
	data := map[string]string{
		"error": err.Error(),
	}

	return s.respond(w, code, data)
}
