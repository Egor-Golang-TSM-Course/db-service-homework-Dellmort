package postgresql

import (
	"github.com/jackc/pgx/v5"
)

type PostrgeSQLStorage struct {
	db *pgx.Conn
}

func NewPostgreSQLStorage(db *pgx.Conn) *PostrgeSQLStorage {
	return &PostrgeSQLStorage{
		db: db,
	}
}
