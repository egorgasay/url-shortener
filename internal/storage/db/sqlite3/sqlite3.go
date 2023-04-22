package sqlite3

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/db/basic"
	"url-shortener/internal/storage/db/service"

	_ "github.com/mattn/go-sqlite3"
)

var (
	_ storage.IStorage = (*Sqlite3)(nil)
)

// Sqlite3 struct with *sql.DB instance.
// It has methods for working with URLs.
type Sqlite3 struct {
	basic.DB
}

// New Sqlite3 struct constructor.
func New(db *sql.DB, path string) service.IRealStorage {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"sqlite", driver)
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

	return &Sqlite3{DB: basic.DB{DB: db}}
}
