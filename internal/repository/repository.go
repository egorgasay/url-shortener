package repository

import (
	"database/sql"
	"url-shortener/internal/storage"
	dbStorage "url-shortener/internal/storage/db"
	filestorage "url-shortener/internal/storage/file"
	mapStorage "url-shortener/internal/storage/map"
)

type Config struct {
	DriverName     storage.Type
	DataSourceName string
}

func New(cfg *Config) (storage.IStorage, error) {
	if cfg == nil {
		panic("конфигурация задана некорректно")
	}

	switch cfg.DriverName {
	case "sqlite3":
		exists := storage.IsDatabaseExist(cfg.DataSourceName)

		db, err := sql.Open(string(cfg.DriverName), cfg.DataSourceName)
		if err != nil {
			return nil, err
		}

		if !exists {
			err = storage.InitDatabase(db)
			if err != nil {
				return nil, err
			}
		}

		return dbStorage.NewRealStorage(db), nil
	case "file":
		filename := cfg.DataSourceName
		return filestorage.NewFileStorage(filename), nil
	default:
		db := mapStorage.NewMapStorage()
		return db, nil
	}
}
