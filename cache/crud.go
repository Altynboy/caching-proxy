package cache

import (
	"bytes"
	"net/http"
	"time"
)

func (c *Cache) Set(url string, resp http.Response, ttl time.Duration) error {
	var buf bytes.Buffer
	resp.Write(&buf)
	return c.set(url, buf.String(), ttl)
}

func (c *Cache) set(key, value string, ttl time.Duration) error {
	expiration := time.Now().Add(ttl).Unix()

	_, err := c.db.Exec(`
		INSERT OR REPLACE INTO cache (key, value, expiration)
		VALUES (?, ?, ?)
	`, key, value, expiration)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Get(key string) (string, bool) {
	var value string
	var expiration int64

	err := c.db.QueryRow(`
		SELECT value, expiration
		FROM cache
		WHERE
			key = ? 
			AND ( expiration > ? OR expiration = 0)
	`, key, time.Now().Unix()).Scan(&value, &expiration)

	if err != nil {
		return "", false
	}

	return value, true
}

func (c *Cache) DeleteAll() error {
	_, err := c.db.Exec("DELETE FROM cache")
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Delete(key string) error {
	_, err := c.db.Exec("DELETE FROM cache WHERE key = ?", key)
	if err != nil {
		return err
	}

	return nil
}
