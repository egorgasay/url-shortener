package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/egorgasay/dockerdb"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"log"
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
	case "sqlite3", "test":
		db, err := sql.Open("sqlite3", cfg.DataSourceCred)
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
		sqlitedb, err := upSqlite(cfg, "file://internal/repository/migrations")
		if err != nil {
			return nil, err
		}

		stmt, err := sqlitedb.Prepare("SELECT id, connectionString FROM DockerDBs WHERE name = ?")
		if err != nil {
			return nil, err
		}

		err = stmt.QueryRow(cfg.Name).Scan(&cfg.VDB.ID, &cfg.DataSourceCred)
		if errors.Is(err, sql.ErrNoRows) {
			stmt, err := sqlitedb.Prepare("INSERT INTO DockerDBs VALUES (?, ?, ?)")
			if err != nil {
				return nil, err
			}

			_, err = stmt.Exec(cfg.Name, cfg.VDB.ID, cfg.DataSourceCred)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}

		ctx := context.TODO()
		err = cfg.VDB.Run(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to run docker storage %w", err)
		}

		sqlitedb.Close()

		return dbStorage.NewRealStorage(cfg.VDB.DB, cfg.DriverName), nil
	case "file":
		filename := cfg.DataSourcePath
		return filestorage.NewFileStorage(filename)
	default:
		db := mapStorage.NewMapStorage()
		return db, nil
	}
}

func upSqlite(cfg *Config, path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DataSourcePath)
	if err != nil {
		return nil, err
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"url-shortener", driver)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = m.Up()
	if err != nil {
		if err.Error() != "no change" {
			log.Fatal(err)
		}
	}

	return db, nil
}
