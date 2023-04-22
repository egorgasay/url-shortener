package basic

import (
	"context"
	"database/sql"
	"fmt"

	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/db/queries"
)

// DB is a basic implementation of the storage.Repository interface.
type DB struct {
	*sql.DB
}

// Ping checks connection with the repository.
func (db *DB) Ping() error {
	if err := db.PingContext(context.Background()); err != nil {
		return fmt.Errorf("error pinging db: %w", err)
	}

	return nil
}

// Shutdown closes the database connection.
func (db *DB) Shutdown() error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("error closing db: %w", err)
	}

	return nil
}

// MarkAsDeleted finds a URL and marks it as deleted.
func (db *DB) MarkAsDeleted(shortURL, cookie string) error {
	stmt, err := queries.GetPreparedStatement(queries.MarkAsDeleted)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}

	_, err = stmt.Exec(
		sql.Named("short", shortURL).Value,
		sql.Named("cookie", cookie).Value,
	)

	if err != nil {
		return fmt.Errorf("error marking as deleted: %w", err)
	}

	return nil
}

// FindMaxID gets len of the repository.
func (db *DB) FindMaxID() (int, error) {
	var id int

	stmt, err := queries.GetPreparedStatement(queries.FindMaxURL)
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %w", err)
	}

	if stmt.QueryRow().Scan(&id) != nil {
		return 0, fmt.Errorf("error finding max id: %w", err)
	}

	return id, nil
}

// GetAllLinksByCookie gets all links ([]schema.URL) by cookie.
func (db *DB) GetAllLinksByCookie(cookie, baseURL string) ([]schema.URL, error) {
	stmt, err := queries.GetPreparedStatement(queries.GetAllLinksByCookie)
	if err != nil {
		return nil, nil
	}

	stm, err := stmt.Query(sql.Named("cookie", cookie).Value)
	if err != nil {
		return nil, fmt.Errorf("error getting links by cookie: %w", err)
	}

	err = stm.Err()
	if err != nil {
		return nil, fmt.Errorf("error getting links by cookie: %w", err)
	}

	var links []schema.URL

	for stm.Next() {
		short, long := "", ""

		err = stm.Scan(&short, &long)
		if err != nil {
			return nil, fmt.Errorf("error getting links by cookie: %w", err)
		}

		links = append(links, schema.URL{LongURL: long, ShortURL: baseURL + short})
	}

	return links, nil
}

// GetLongLink gets a long link from the repository.
func (db *DB) GetLongLink(shortURL string) (longURL string, err error) {
	stmt, err := queries.GetPreparedStatement(queries.GetLongLink)
	if err != nil {
		return "", fmt.Errorf("error preparing statement: %w", err)
	}

	var isDeleted = false
	if stmt.QueryRow(sql.Named("short", shortURL).Value).Scan(&longURL, &isDeleted) != nil {
		return "", fmt.Errorf("error getting long link: %w", err)
	}

	if isDeleted {
		return "", fmt.Errorf("error getting long link: %w", storage.ErrDeleted)
	}

	return longURL, nil
}

// AddLink adds a link to the repository.
func (db *DB) AddLink(longURL, shortURL, cookie string) (string, error) {
	stmt, err := queries.GetPreparedStatement(queries.InsertURL)
	if err != nil {
		return "", fmt.Errorf("error preparing statement: %w", err)
	}

	_, err = stmt.Exec(
		sql.Named("long", longURL).Value,
		sql.Named("short", shortURL).Value,
		sql.Named("cookie", cookie).Value,
	)

	if err != nil {
		return "", fmt.Errorf("error adding link: %w", err)
	}

	return shortURL, nil
}
