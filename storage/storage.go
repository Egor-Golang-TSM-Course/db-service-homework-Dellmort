package storage

import (
	"context"
	"homework/internal/models"
	"time"
)

// Storage ...
type Storage interface {
	User
	Post
	Comment
	Tags
}

type User interface {
	CreateUser(ctx context.Context, user *models.User) (int, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
}

type Post interface {
	CreatePost(ctx context.Context, post *models.Post) (int, error)
	GetPostByID(ctx context.Context, id int) (*models.Post, error)
	GetPosts(ctx context.Context, tags []string, date *time.Time) ([]*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, id int) error
}

type Comment interface {
	CreateComment(ctx context.Context, comment *models.Comment) (int, error)
	GetCommentsByPostID(ctx context.Context, postID int) ([]*models.Comment, error)
}

type Tags interface {
	CreateTags(context.Context, *models.Tags) error
	GetTagsByPostID(ctx context.Context, postID int) ([]string, error)
}
