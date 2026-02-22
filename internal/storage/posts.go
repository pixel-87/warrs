package storage

import (
	"errors"
	"fmt"

	"github.com/pixel-87/warss/internal/models"
)

func (d *DB) AddPosts (feedID int, posts []models.Post) error {
	query := `INSERT INTO posts (
		feed_id,
		title,
		link,
		content,
		published_at,
		updated_at,
	)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(link) DO NOTHING;`
	
	for i := range posts {
		_, err := d.conn.Exec(
			query,
			feedID,
			posts[i].Title,
			posts[i].Link,
			posts[i].Content
			posts[i].PublishedAt
			posts[i].UpdatedAt
		)

		if err != nil {
			return fmt.Errorf("failed to insert post %q for feed %d: %w", posts[i].Title, feedID, err)
		}
		return nil
	}
}
