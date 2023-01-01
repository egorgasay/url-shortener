package repository

import (
	"database/sql"
	"os"
	"url-shortener/internal/storage"
	dbStorage "url-shortener/internal/storage/db"
	filestorage "url-shortener/internal/storage/file"
	mapStorage "url-shortener/internal/storage/map"
)

type Config struct {
	DriverName     string
	DataSourceName string
}

type Repository struct {
	repo storage.IStorage
}

func New(cfg *Config) (*Repository, error) {
	if cfg == nil {
		panic("конфигурация задана некорректно")
	}

	switch cfg.DriverName {
	case "sqlite3":
		db, err := sql.Open(cfg.DriverName, cfg.DataSourceName)
		if err != nil {
			return nil, err
		}

		return &Repository{repo: dbStorage.NewRealStorage(db)}, nil
	case "file":
		filename := cfg.DataSourceName
		if name, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
			filename = name
		}

		return &Repository{repo: filestorage.NewFileStorage(filename)}, nil

	default:
		db := mapStorage.NewMapStorage()
		return &Repository{repo: &db}, nil
	}
}
