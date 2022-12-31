package repository

import (
	"database/sql"
	"url-shortener/internal/storage"
	dbStorage "url-shortener/internal/storage/db"
	mapStorage "url-shortener/internal/storage/map"
)

type Config struct {
	DriverName     string
	DataSourceName string
}

func New(cfg *Config) (*storage.Storage, error) {
	if cfg == nil {
		panic("конфигурация задана некорректно")
	}

	if cfg.DriverName == "map" {
		db := mapStorage.NewMapStorage()
		return &storage.Storage{
			DB: &db,
		}, nil
	}

	db, err := sql.Open(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		return nil, err
	}

	return &storage.Storage{
		DB: dbStorage.NewRealStorage(db),
	}, nil
}
