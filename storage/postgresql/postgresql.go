package postgresql

import (
	"context"
	"homework/internal/models"

	"github.com/jmoiron/sqlx"
)

type PostrgeSQLStorage struct {
	db *sqlx.DB
}

func NewPostgreSQLStorage(db *sqlx.DB) *PostrgeSQLStorage {
	return &PostrgeSQLStorage{
		db: db,
	}
}

func (p *PostrgeSQLStorage) CreateUser(ctx context.Context, user *models.User) (int, error) {
	query := "INSERT INTO users (name, password, login, created) VALUES ($1, $2, $3, $4) RETURNING id"

	row := p.db.QueryRowContext(ctx, query, user.Name, user.Password, user.Login, user.Created)
	if row.Err() != nil {
		return 0, row.Err()
	}
	var id int
	row.Scan(&id)

	return id, nil
}

func (p *PostrgeSQLStorage) GetUser(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE login = $1"

	row := p.db.QueryRowxContext(ctx, query, login)
	if err := row.StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
