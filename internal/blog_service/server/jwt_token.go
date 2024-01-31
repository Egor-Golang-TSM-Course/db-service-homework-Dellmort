package server

import (
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var EXPDATE = time.Now().Add(72 * time.Hour).Unix()

type TokenManager interface {
	GenerateJWTToken(id int, login string) (string, error)
	ParseJWTtoken(jwtToken string) (jwt.MapClaims, error)
}

type JWTManager struct {
	jwtSecretKey string
}

func NewJWTManager(jwtSecretKey string) *JWTManager {
	return &JWTManager{
		jwtSecretKey: jwtSecretKey,
	}
}

func (m *JWTManager) GenerateJWTToken(id int, login string) (string, error) {
	const f = "server.JWTGenerate"

	payload := jwt.MapClaims{
		"id":  id,
		"exp": EXPDATE,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString([]byte(m.jwtSecretKey))
	if err != nil {
		slog.Error(f, slog.String("ERR", err.Error()))
		return "", err
	}
	return t, nil
}

func (m *JWTManager) ParseJWTtoken(jwtToken string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(m.jwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}
