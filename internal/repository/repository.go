package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/egorgasay/dockerdb"
	"url-shortener/internal/storage"
	dbStorage "url-shortener/internal/storage/db"
	filestorage "url-shortener/internal/storage/file"
	mapStorage "url-shortener/internal/storage/map"
)

type Config struct {
	DriverName     storage.Type
	DataSourceCred string
	DataSourcePath string
	VDB            *dockerdb.VDB
	Name           string
}

func New(cfg *Config) (storage.IStorage, error) {
	if cfg == nil {
		panic("конфигурация задана некорректно")
	}

	switch cfg.DriverName {
	case "sqlite3":
		db, err := upSqlite(cfg, "sqlite3-schema.sql")
		if err != nil {
			return nil, err
		}

		return dbStorage.NewRealStorage(db, cfg.DriverName), nil
	case "mysql", "postgres":
		var db *sql.DB
		var err error

		if cfg.VDB == nil {
			db, err = sql.Open(string(cfg.DriverName), cfg.DataSourceCred)
			if err != nil {
				return nil, err
			}
			return dbStorage.NewRealStorage(db, cfg.DriverName), nil
		}

		cfg.DataSourcePath = "dockerDBs"
		sqlitedb, err := upSqlite(cfg, "DockerDBs-schema.sql")
		if err != nil {
			return nil, err
		}

		stmt, err := sqlitedb.Prepare("SELECT id, connectionString FROM DockerDBs WHERE name = ?")
		if err != nil {
			return nil, err
		}

		row := stmt.QueryRow(cfg.Name)

		err = row.Err()
		if err != nil {
			return nil, err
		}

		err = row.Scan(&cfg.VDB.ID, &cfg.DataSourceCred)
		if err != sql.ErrNoRows && err != nil {
			return nil, err
		}

		ctx := context.TODO()
		err = cfg.VDB.Run(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to run docker storage %w", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			stmt, err := sqlitedb.Prepare("INSERT INTO DockerDBs VALUES (?, ?, ?)")
			if err != nil {
				return nil, err
			}

			_, err = stmt.Exec(cfg.Name, cfg.VDB.ID, cfg.DataSourceCred)
			if err != nil {
				return nil, err
			}
		}
		sqlitedb.Close()

		return dbStorage.NewRealStorage(cfg.VDB.DB, cfg.DriverName), nil
	case "file":
		filename := cfg.DataSourcePath
		return filestorage.NewFileStorage(filename), nil
	default:
		db := mapStorage.NewMapStorage()
		return db, nil
	}
}

func upSqlite(cfg *Config, schema string) (*sql.DB, error) {
	exists := storage.IsDBSqliteExist(cfg.DataSourcePath)

	db, err := sql.Open("sqlite3", cfg.DataSourcePath)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = storage.InitDatabase(db, schema)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
