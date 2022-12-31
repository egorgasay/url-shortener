package mapstorage

import "errors"

func (s MapStorage) GetLongLink(shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	longURL, ok := s.container[shortURL]
	if !ok {
		return longURL, errors.New("короткой ссылки не существует")
	}

	return longURL, nil
}
