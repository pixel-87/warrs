package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

// init sqlite3 file and create tables if they don't exist
func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS feeds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT UNIQUE NOT NULL,
		title TEXT
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("error creating feed table: %w", err)
	}

	query = `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		feed_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		link TEXT UNIQUE NOT NULL,
		content TEXT,
		published_at DATETIME,
		read BOOLEAN DEFAULT 0,
		FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("error creating posts table: %w", err)
	}

	return &DB{conn: db}, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}
