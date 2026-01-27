package storage

import (
	"github.com/pixel-87/warss/internal/models"
)

func (d *DB) AddFeed(url, title string) error {
	query := `INSERT INTO feeds (url, title) VALUES (?, ?)`
	_, err := d.conn.Exec(query, url, title)
	return err
}

func (d *DB) GetFeeds() ([]models.Feed, error) {
	rows, err := d.conn.Query("SELECT id, url, title FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []models.Feed
	for rows.Next() {
		var f models.Feed
		if err := rows.Scan(&f.ID, &f.URL, &f.Title); err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}
	return feeds, nil
}
