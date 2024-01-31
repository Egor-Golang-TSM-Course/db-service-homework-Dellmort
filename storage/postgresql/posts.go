package postgresql

import (
	"context"
	"errors"
	"homework/internal/models"
	"strings"
	"time"
)

var (
	ErrNoSearchResult = errors.New("no rows this tags and date")
)

// TODO: Изменить порядок сканирования как в базе данных
func (p *PostrgeSQLStorage) CreatePost(ctx context.Context, post *models.Post) (int, error) {
	query := "INSERT INTO posts (name, owner_id, message, created) VALUES ($1, $2, $3, $4) RETURNING id"
	row := p.db.QueryRow(ctx, query, post.Name, post.OwnerID, post.Message, post.Created)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// TODO: FiX
func (p *PostrgeSQLStorage) GetPosts(ctx context.Context, tags []string, date *time.Time) ([]*models.Post, error) {
	query := `
	SELECT * FROM posts
	LEFT JOIN tags ON tags.post_id = posts.id
	WHERE tags IS NOT NULL AND post_id IS NOT NULL
`

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)
	for rows.Next() {
		post := models.NewPost()
		err = rows.Scan(
			&post.ID,
			&post.Name,
			&post.Message,
			&post.Created,
			&post.OwnerID,
			&post.Tags.Tags,
			&post.Tags.PostID,
		)
		if err != nil {
			return nil, err
		}

		if (tags != nil || len(tags) > 0) && (date != nil && !date.IsZero()) {
			var count int
			tagsCap := strings.Fields(post.Tags.Tags)
			for _, tag := range tagsCap {
				for _, expTag := range tags {
					if tag == expTag {
						count++
					}
				}
				if count == len(tags) {
					break
				}

				return nil, ErrNoSearchResult
			}
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *PostrgeSQLStorage) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	query := `
	SELECT * FROM posts WHERE posts.id=$1`

	row := p.db.QueryRow(ctx, query, id)
	post := models.NewPost()
	err := row.Scan(&post.ID,
		&post.Name,
		&post.Message,
		&post.Created,
		&post.OwnerID,
	)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *PostrgeSQLStorage) UpdatePost(ctx context.Context, post *models.Post) error {
	query := "UPDATE posts SET name = $1, message = $2 WHERE id = $3"

	_, err := p.db.Exec(ctx, query, post.Name, post.Message, post.ID)
	return err
}

func (p *PostrgeSQLStorage) DeletePost(ctx context.Context, id int) error {
	query := "DELETE FROM posts WHERE id = $1"

	_, err := p.db.Exec(ctx, query, id)
	return err
}
