package dbStorage

import "database/sql"

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
