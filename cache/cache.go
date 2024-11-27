package cache

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Cache struct {
	db *sql.DB
}

func New(path string) (*Cache, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Can't open db file")
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cache (
			key TEXT PRIMARY KEY,
			value TEXT,
			expiration INTEGER
		);
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Cache{db: db}, err
}

func (c *Cache) Close() error {
	return c.db.Close()
}
