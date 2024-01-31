package postgresql

import (
	"context"
	"homework/internal/models"
)

func (p *PostrgeSQLStorage) CreateTags(ctx context.Context, tags *models.Tags) error {
	query := "INSERT INTO tags (tags, post_id) VALUES ($1, $2)"

	_, err := p.db.Exec(ctx, query, tags.Tags, tags.PostID)
	return err
}

func (p *PostrgeSQLStorage) GetTagsByPostID(ctx context.Context, postID int) ([]string, error) {
	var tags []string

	query := "SELECT tags FROM tags WHERE post_id = $1"
	row, err := p.db.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var tag string
		if err := row.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
