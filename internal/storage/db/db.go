package dbstorage

import (
	"database/sql"
	"strings"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
	shortenalgorithm "url-shortener/pkg/shortenAlgorithm"
)

type RealStorage struct {
	DB *sql.DB
}

const DBStorageType storage.Type = "sqlite3"

func NewRealStorage(db *sql.DB) storage.IStorage {
	return &RealStorage{DB: db}
}

func (s RealStorage) AddLink(longURL string, id int, cookie string) (string, error) {
	shortURL, err := shortenalgorithm.GetShortName(id)
	if err != nil {
		return "", err
	}

	stmt := "INSERT INTO urls (long, short, cookie) VALUES (?, ?, ?)"

	_, err = s.DB.Exec(stmt, longURL, shortURL, cookie)

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
	stm := s.DB.QueryRow("SELECT long FROM urls WHERE short = ?", shortURL)
	err = stm.Scan(&longURL)

	return longURL, err
}

func (s RealStorage) GetAllLinksByCookie(cookie string) ([]schema.URL, error) {
	stm, err := s.DB.Query("SELECT short, long FROM urls WHERE cookie = ?", cookie)
	if err != nil {
		return nil, err
	}
	err = stm.Err()
	if err != nil {
		return nil, err
	}

	var URLs []schema.URL
	for stm.Next() {
		tmp := ""

		err = stm.Scan(&tmp)
		if err != nil {
			return nil, err
		}

		lineArr := strings.Split(tmp, " ")

		URLs = append(URLs, schema.URL{LongURL: lineArr[1], ShortURL: lineArr[2]})
	}

	return URLs, err
}
