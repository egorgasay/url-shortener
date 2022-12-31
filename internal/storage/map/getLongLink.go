package mapstorage

import "errors"

func (s MapStorage) GetLongLink(shortURL string) (string, error) {
	longURL, ok := s.container[shortURL]
	if !ok {
		return longURL, errors.New("короткой ссылки не существует")
	}

	return longURL, nil
}
