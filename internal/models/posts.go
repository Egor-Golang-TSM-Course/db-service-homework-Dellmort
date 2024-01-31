package models

import (
	"time"
)

type Post struct {
	ID      int       `db:"id" json:"id,omitempty"`
	Name    string    `db:"name" validate:"required" json:"name,omitempty"`
	Message string    `db:"message" validate:"required" json:"message,omitempty"`
	Created time.Time `db:"created" json:"created,omitempty"`
	OwnerID int       `db:"owner_id" json:"owner_id,omitempty"`
	*Comment
	*Tags
}

type Comment struct {
	ID      int       `db:"id" json:"id,omitempty"`
	PostID  int       `db:"post_id" json:"post_id,omitempty"`
	OwnerID int       `db:"owner_id" json:"owner_id,omitempty"`
	Message string    `db:"message" json:"message,omitempty"`
	Created time.Time `db:"created" json:"created,omitempty"`
}

type Tags struct {
	Tags   string `db:"tags" json:"tags,omitempty"`
	PostID int    `db:"post_id" json:"post_id,omitempty"`
}

func NewPost() *Post {
	return &Post{
		Comment: new(Comment),
		Tags:    new(Tags),
	}
}
