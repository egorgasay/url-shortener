package storage

import (
	"database/sql"
)

type IRealStorage interface {
	Ping() error
	Close() error
	Query(string, ...any) (*sql.Rows, error)
	Exec(string, ...any) (sql.Result, error)
	QueryRow(string, ...any) *sql.Row
}

type RealStorage struct {
	IRealStorage
}

func NewRealStorage(db IRealStorage) RealStorage {
	return RealStorage{db}
}

type MapStorage map[string]string

func NewMapStorage() MapStorage {
	db := make(MapStorage, 10)
	return db
}
