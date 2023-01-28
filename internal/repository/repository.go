package repository

import (
	"database/sql"
	"errors"
	"url-shortener/internal/dockerdb"
	"url-shortener/internal/storage"
	dbStorage "url-shortener/internal/storage/db"
	filestorage "url-shortener/internal/storage/file"
	mapStorage "url-shortener/internal/storage/map"
)

type Config struct {
	DriverName     storage.Type
	DataSourceCred string
	DataSourcePath string
	DockerDB       *dockerdb.DockerDB
}

func New(cfg *Config) (storage.IStorage, error) {
	if cfg == nil {
		panic("конфигурация задана некорректно")
	}

	switch cfg.DriverName {
	case "sqlite3":
		db, err := upSqlite(cfg, "schema.sql")
		if err != nil {
			return nil, err
		}

		return dbStorage.NewRealStorage(db), nil
	case "mysql", "postgres":
		var db *sql.DB
		var err error

		if cfg.DockerDB != nil {
			cfg.DataSourcePath = "dockerDBs"
			sqlitedb, err := upSqlite(cfg, "sqlite-schema.sql")
			if err != nil {
				return nil, err
			}

			stmt, err := sqlitedb.Prepare("SELECT id, connectionString FROM DockerDBs WHERE name = ?")
			if err != nil {
				return nil, err
			}

			row := stmt.QueryRow(cfg.DockerDB.Conf.DB.Name)

			err = row.Err()
			if err != nil {
				return nil, err
			}

			err = row.Scan(&cfg.DockerDB.ID, &cfg.DataSourceCred)
			if err != sql.ErrNoRows && err != nil {
				return nil, err
			}

			if cfg.DataSourceCred == "" {
				db, cfg.DataSourceCred = cfg.DockerDB.Setup("")
			} else {
				db, _ = cfg.DockerDB.Setup(cfg.DataSourceCred)
			}

			if errors.Is(err, sql.ErrNoRows) {
				stmt, err := sqlitedb.Prepare("INSERT INTO DockerDBs VALUES (?, ?, ?)")
				if err != nil {
					return nil, err
				}

				_, err = stmt.Exec(cfg.DockerDB.Conf.DB.Name, cfg.DockerDB.ID, cfg.DataSourceCred)
				if err != nil {
					return nil, err
				}
			}
			sqlitedb.Close()
		} else {
			db, err = sql.Open(string(cfg.DriverName), cfg.DataSourceCred)
			if err != nil {
				return nil, err
			}
		}

		used := storage.IsDBUsedBefore(db)

		if !used {
			err := storage.InitDatabase(db, "schema.sql")
			if err != nil {
				return nil, err
			}
		}

		return dbStorage.NewRealStorage(db), nil
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
