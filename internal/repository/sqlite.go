package repository

import (
	"database/sql"
	"log"
)

type Config struct {
	DriverName     string
	DataSourceName string
}

func NewSqliteDb(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
