package storage

import (
	"fmt"

	"github.com/pixel-87/warss/internal/models"
)

func (d *DB) AddFeed(url, title string) error {
	query := `INSERT INTO feeds (url, title) VALUES (?, ?)`
	_, err := d.conn.Exec(query, url, title)
	if err != nil {
		return fmt.Errorf("failed to add feed %q: %w", url, err)
	}
	return nil
}

func (d *DB) GetFeeds() ([]models.Feed, error) {
	rows, err := d.conn.Query("SELECT id, url, title FROM feeds")
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			fmt.Printf("warning: rows.Close error: %v\n", cerr)
		}
	}()

	var feeds []models.Feed
	for rows.Next() {
		var f models.Feed
		if err := rows.Scan(&f.ID, &f.URL, &f.Title); err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating feeds: %w", err)
	}
	return feeds, nil
}

func (d *DB) UpdateFeed(f models.Feed) error {
	query := `
		UPDATE feeds
		SET url = ?, title = ?
		WHERE id = ?
	`

	_, err := d.conn.Exec(query, f.URL, f.Title, f.ID)
	if err != nil {
		return fmt.Errorf("failed to update feed %d: %w", f.ID, err)
	}
	return nil
}

func (d *DB) DeleteFeed(id int) error {
	query := `DELETE FROM feeds WHERE id = ?`
	_, err := d.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not delete feed: %w", err)
	}
	return nil
}
