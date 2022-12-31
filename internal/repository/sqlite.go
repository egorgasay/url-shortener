package repository

import (
	"database/sql"
	"url-shortener/internal/storage"
)

type Config struct {
	DriverName     string
	DataSourceName string
}

//type IStorage interface {
//	Ping() error
//	Close() error
//	Exec(string, ...any) (sql.Result, error)
//	QueryRow(string, ...any) *sql.Row
//}

type IStorage interface {
	FindMaxID() (int, error)
	AddLink(longURL string, id int) (string, error)
	GetLongLink(shortURL string) (longURL string, err error)
}

type Storage struct {
	DB IStorage
}

func New(cfg *Config) (*Storage, error) {
	if cfg == nil {
		panic("конфигурация задана некорректно")
	}

	if cfg.DriverName == "map" {
		db := storage.NewMapStorage()
		return &Storage{
			DB: db,
		}, nil
	}

	db, err := sql.Open(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		return nil, err
	}

	return &Storage{
		DB: storage.NewRealStorage(db),
	}, nil
}
