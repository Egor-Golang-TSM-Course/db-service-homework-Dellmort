package postgresql

import (
	"context"
	"homework/internal/models"
	"log/slog"
)

func (p *PostrgeSQLStorage) CreateComment(ctx context.Context, comment *models.Comment) (int, error) {
	query := "INSERT INTO comments (post_id, message, owner_id, created) VALUES ($1, $2, $3, $4) RETURNING id"

	row := p.db.QueryRow(ctx, query, comment.PostID, comment.Message, comment.OwnerID, comment.Created)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *PostrgeSQLStorage) GetCommentsByPostID(ctx context.Context, postID int) ([]*models.Comment, error) {
	const f = "postgresql.GetCommentsByPostID"

	query := "SELECT * FROM comments WHERE post_id=$1"

	rows, err := p.db.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	comments := make([]*models.Comment, 0)
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID,
			&comment.Message,
			&comment.Created,
			&comment.PostID,
			&comment.OwnerID,
		)
		if err != nil {
			slog.Error(f, slog.String("ERR", err.Error()))
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}
