package dbstorage

import (
	"database/sql"
	"log"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/db/mysql"
	"url-shortener/internal/storage/db/postgres"
	"url-shortener/internal/storage/db/queries"
	"url-shortener/internal/storage/db/service"
	"url-shortener/internal/storage/db/sqlite3"
)

// DBStorageType postgres type.
const DBStorageType storage.Type = "postgres"

// NewRealStorage constructor for storage.IStorage with db implementation.
func NewRealStorage(db *sql.DB, vendor storage.Type) service.IRealStorage {
	var irs service.IRealStorage
	switch vendor {
	case "postgres":
		irs = postgres.New(db, "file://migrations/postgres")
	case "mysql":
		irs = mysql.New(db, "file://migrations/mysql")
	case "sqlite3":
		irs = sqlite3.New(db, "file://migrations/sqlite3")
	case "test":
		irs = sqlite3.New(db, "file://../../migrations/sqlite3")
	}

	var err error
	if vendor != "test" {
		err = queries.Prepare(db, string(vendor))
	} else {
		err = queries.Prepare(db, "sqlite3")
	}

	if err != nil {
		log.Fatal("failed to prepare queries: ", err)
	}

	return irs
}
