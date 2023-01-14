package dbstorage

import (
	"database/sql"
	"url-shortener/internal/storage"
	shortenalgorithm "url-shortener/pkg/shortenAlgorithm"
)

type RealStorage struct {
	DB *sql.DB
}

func NewRealStorage(db *sql.DB) storage.IStorage {
	return &RealStorage{DB: db}
}

func (s RealStorage) AddLink(longURL string, id int) (string, error) {
	shortURL := shortenalgorithm.GetShortName(id)
	stmt := "INSERT INTO urls (long, short) VALUES (?, ?)"

	_, err := s.DB.Exec(stmt, longURL, shortURL)

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
