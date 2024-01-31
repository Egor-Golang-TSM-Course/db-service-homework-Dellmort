package postgresql

import (
	"context"
	"homework/internal/models"
)

func (p *PostrgeSQLStorage) CreateUser(ctx context.Context, user *models.User) (int, error) {
	query := "INSERT INTO users (name, password, login, created) VALUES ($1, $2, $3, $4) RETURNING id"

	row := p.db.QueryRow(ctx, query, user.Name, user.Password, user.Login, user.Created)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *PostrgeSQLStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE login = $1"

	row := p.db.QueryRow(ctx, query, login)
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Login,
		&user.Created,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostrgeSQLStorage) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE id = $1"

	row := p.db.QueryRow(ctx, query, id)
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Login,
		&user.Created,
	); err != nil {
		return nil, err
	}

	return &user, nil
}
