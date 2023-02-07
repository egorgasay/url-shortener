package dbstorage

import (
	"database/sql"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/db/mysql"
	"url-shortener/internal/storage/db/postgres"
	"url-shortener/internal/storage/db/service"
	"url-shortener/internal/storage/db/sqlite3"
)

const DBStorageType storage.Type = "postgres"

func NewRealStorage(db *sql.DB, vendor storage.Type) service.IRealStorage {
	switch vendor {
	case "postgres":
		return postgres.New(db)
	case "mysql":
		return mysql.New(db)
	case "sqlite3":
		return sqlite3.New(db)
	}

	return nil
}
