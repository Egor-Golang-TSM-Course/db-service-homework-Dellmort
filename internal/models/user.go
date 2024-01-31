package models

import "time"

type User struct {
	ID       int       `db:"id" json:"id"`
	Name     string    `db:"name" json:"name" validate:"required"`
	Login    string    `db:"login" json:"login" validate:"required"`
	Password string    `db:"password" json:"password" validate:"required"`
	Created  time.Time `db:"created" json:"created"`
}
