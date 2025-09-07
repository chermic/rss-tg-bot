package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

type PostLink struct {
	ID        int64
	Link      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("sqlite3", "./db/rss-tg-bot.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS rss_posts (id INTEGER PRIMARY KEY AUTOINCREMENT, url TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME DEFAULT CURRENT_TIMESTAMP)")

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) AddLink(link string) error {
	_, err := d.db.Exec("INSERT INTO rss_posts (url) VALUES (?)", link)
	return err
}

func (d *Database) HasLink(link string) (bool, error) {
	rows, err := d.db.Query("SELECT COUNT(*) FROM rss_posts WHERE url = ?", link)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	rows.Next()

	var count int
	if err := rows.Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}
