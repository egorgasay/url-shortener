package storage

import "errors"

func (s MapStorage) GetLongLink(shortURL string) (string, error) {
	longURL, ok := s[shortURL]
	if !ok {
		return longURL, errors.New("короткой ссылки не существует")
	}

	return longURL, nil
}

func (s RealStorage) GetLongLink(shortURL string) (longURL string, err error) {
	stm := s.QueryRow("SELECT long FROM urls WHERE short = ?", shortURL)
	err = stm.Scan(&longURL)

	return longURL, err
}
