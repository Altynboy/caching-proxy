package cache

import (
	"bytes"
	"database/sql"
	"io"
	"log"
	"net/http"
	"time"
)

func (c *Cache) Set(url string, resp http.Response, ttl time.Duration) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	var buf bytes.Buffer
	err = resp.Write(&buf)
	if err != nil {
		return err
	}
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

func (c *Cache) Get(key string) (string, bool, error) {
	var value string
	var expiration int64
	now := time.Now().Unix()
	err := c.db.QueryRow(`
		SELECT value, expiration
		FROM cache
		WHERE
			key = ? 
	    	AND ( expiration > ? OR expiration = 0)
	`, key, now).Scan(&value, &expiration)
	if err == sql.ErrNoRows {
		return "", false, nil
	} else if err != nil {
		log.Printf("url: %q\n", key)
		log.Printf("now: %d\n", now)
		log.Fatalf("Error while getting cache %s", err)
		return "", false, err
	}

	return value, true, nil
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
