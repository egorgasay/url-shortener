package dbstorage

import "url-shortener/internal/storage/shortenAlgorithm"

func (s RealStorage) AddLink(longURL string, id int) (string, error) {
	shortURL := shortenalgorithm.GetShortName(id)
	stmt := "INSERT INTO urls (long, short) VALUES (?, ?)"

	_, err := s.DB.Exec(stmt, longURL, shortURL)

	if err != nil {
		return "", err
	}

	return shortURL, nil
}
