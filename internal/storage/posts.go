package storage

import (
	"context"
	"fmt"

	"github.com/pixel-87/warss/internal/models"
)

func (d *DB) AddPosts(feedID int, posts []models.Post) error {
	query := `INSERT INTO posts (
		feed_id,
		title,
		link,
		content,
		published_at,
		updated_at
	)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(link) DO NOTHING;`

	for i := range posts {
		_, err := d.conn.Exec(
			query,
			feedID,
			posts[i].Title,
			posts[i].Link,
			posts[i].Content,
			posts[i].PublishedAt,
			posts[i].UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to insert post %q for feed %d: %w", posts[i].Title, feedID, err)
		}
	}
	return nil
}

func (d *DB) GetPost(ctx context.Context, postID int) (models.Post, error) {
	query := `
		SELECT feed_id, title, link, content, published_at, updated_at
		FROM posts
		WHERE id = ?;
	`

	var p models.Post
	err := d.conn.QueryRowContext(ctx, query, postID).Scan(
		&p.FeedID,
		&p.Title,
		&p.Link,
		&p.Content,
		&p.PublishedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return models.Post{}, fmt.Errorf("failed to get post id:%d, %w", postID, err)
	}
	return p, nil
}
