package repository

import (
	"database/sql"
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

		exists := IsDatabaseExist(cfg.DataSourceName)

		db, err := sql.Open(cfg.DriverName, cfg.DataSourceName)
		if err != nil {
			return nil, err
		}

		if !exists {
			err = InitDatabase(db)
			if err != nil {
				return nil, err
			}
		}

		return &Repository{repo: dbStorage.NewRealStorage(db)}, nil
	case "file":
		filename := cfg.DataSourceName
		return &Repository{repo: filestorage.NewFileStorage(filename)}, nil

	default:
		db := mapStorage.NewMapStorage()
		return &Repository{repo: db}, nil
	}
}
