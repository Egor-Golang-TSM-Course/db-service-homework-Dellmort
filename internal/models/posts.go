package models

import "time"

type Post struct {
	ID      int        `db:"id"`
	Name    string     `db:"name"`
	Owner   string     `db:"owner_name"`
	Created *time.Time `db:"created"`
	Comment Comment
	Tags    Tags
}

type Comment struct {
	ID      int        `db:"id"`
	PostID  int        `db:"post_id"`
	Message string     `db:"message"`
	Created *time.Time `db:"created"`
}

type Tags struct {
	Tags []string `db:"tags"`
}
