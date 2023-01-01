package mapstorage

import (
	"url-shortener/internal/storage/shortenAlgorithm"
)

func (s *MapStorage) AddLink(longURL string, id int) (string, error) {
	shortURL := shortenAlgorithm.GetShortName(id)

	s.mu.RLock()
	defer s.mu.RUnlock()
	s.container[shortURL] = longURL

	return shortURL, nil
}
