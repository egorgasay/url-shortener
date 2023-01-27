package dbstorage

import (
	"context"
	"database/sql"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type RealStorage struct {
	DB *sql.DB
}

const DBStorageType storage.Type = "postgres"

func NewRealStorage(db *sql.DB) storage.IStorage {
	return &RealStorage{DB: db}
}

func (s RealStorage) AddLink(longURL, shortURL, cookie string) (string, error) {
	stmt := "INSERT INTO urls (long, short, cookie) VALUES ($1, $2, $3)"

	_, err := s.DB.Exec(stmt, longURL, shortURL, cookie)

	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s RealStorage) FindMaxID() (int, error) {
	var id int

	stm := s.DB.QueryRow("SELECT MAX(id) FROM urls")
	err := stm.Scan(&id)

	return id, err
}

func (s RealStorage) GetLongLink(shortURL string) (longURL string, err error) {
	stm := s.DB.QueryRow("SELECT long FROM urls WHERE short = $1", shortURL)
	err = stm.Scan(&longURL)

	return longURL, err
}

func (s RealStorage) GetAllLinksByCookie(cookie, baseURL string) ([]schema.URL, error) {
	stm, err := s.DB.Query("SELECT short, long FROM urls WHERE cookie = $1", cookie)
	if err != nil {
		return nil, err
	}

	err = stm.Err()
	if err != nil {
		return nil, err
	}

	var URLs []schema.URL

	for stm.Next() {
		short, long := "", ""

		err = stm.Scan(&short, &long)
		if err != nil {
			return nil, err
		}

		URLs = append(URLs, schema.URL{LongURL: long, ShortURL: baseURL + short})
	}

	return URLs, err
}

func (s RealStorage) Ping() error {
	ctx := context.TODO()
	return s.DB.PingContext(ctx)
}
