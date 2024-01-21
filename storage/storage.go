package storage

import (
	"context"
	"homework/internal/models"
)

// Storage ...
type Storage interface {
	CreateUser(ctx context.Context, user *models.User) (int, error)
	GetUser(ctx context.Context, login string) (*models.User, error)
}
