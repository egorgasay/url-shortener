package postgres

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"log"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/db/basic"
	"url-shortener/internal/storage/db/queries"
	"url-shortener/internal/storage/db/service"
	shortenalgorithm "url-shortener/pkg/shortenAlgorithm"
)

var (
	_ storage.IStorage = (*Postgres)(nil)
)

// Postgres struct with *sql.DB instance.
// It has methods for working with URLs.
type Postgres struct {
	basic.DB
}

// New Postgres struct constructor.
func New(db *sql.DB, path string) service.IRealStorage {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = m.Up()
	if err != nil {
		if err.Error() != "no change" {
			log.Fatal(err)
		}
	}

	return &Postgres{DB: basic.DB{DB: db}}
}

// AddLink adds a link to the repository.
func (p *Postgres) AddLink(longURL, shortURL, cookie string) (string, error) {
	stmt, err := queries.GetPreparedStatement(queries.InsertURL)
	if err != nil {
		return "", err
	}

	GetShortLinkSTMT, err := queries.GetPreparedStatement(queries.GetShortLink)
	if err != nil {
		return "", err
	}

	_, err = stmt.Exec(
		sql.Named("long", longURL).Value,
		sql.Named("short", shortURL).Value,
		sql.Named("cookie", cookie).Value,
	)

	if err == nil {
		return shortURL, nil
	}

	e, ok := err.(*pq.Error)
	if !ok {
		log.Println("shouldn't be this ", err)
		return "", err
	}

	if e.Code != pgerrcode.UniqueViolation {
		return "", fmt.Errorf("UniqueViolation error: %s", err)
	}

	row := GetShortLinkSTMT.QueryRow(sql.Named("long", longURL).Value)
	if row.Err() != nil {
		return "", err
	}

	err = row.Scan(&shortURL)
	if err == nil {
		return shortURL, service.ErrExists
	}

	log.Println("cycle")

	lastID, err := p.FindMaxID()
	if err != nil {
		return "", err
	}

	shortURL, err = shortenalgorithm.GetShortName(lastID + 1)
	if err != nil {
		return "", err
	}

	return p.AddLink(longURL, shortURL, cookie)
}
