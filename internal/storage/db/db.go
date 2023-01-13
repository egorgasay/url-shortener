package dbstorage

import (
	"database/sql"
	"url-shortener/internal/storage"
)

type RealStorage struct {
	DB *sql.DB
}

func NewRealStorage(db *sql.DB) storage.IStorage {
	return &RealStorage{DB: db}
}
