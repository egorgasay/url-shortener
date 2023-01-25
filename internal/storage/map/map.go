package mapstorage

import (
	"errors"
	"sync"
	"url-shortener/internal/storage"
	shortenalgorithm "url-shortener/pkg/shortenAlgorithm"
)

type MapStorage struct {
	mu        sync.RWMutex
	container map[string]string
}

func NewMapStorage() storage.IStorage {
	db := make(map[string]string, 10)
	return &MapStorage{container: db}
}

func (s *MapStorage) AddLink(longURL string, id int) (string, error) {
	shortURL := shortenalgorithm.GetShortName(id)

	s.mu.RLock()
	defer s.mu.RUnlock()
	s.container[shortURL] = longURL

	return shortURL, nil
}

func (s *MapStorage) FindMaxID() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.container), nil
}

func (s *MapStorage) GetLongLink(shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	longURL, ok := s.container[shortURL]
	if !ok {
		return longURL, errors.New("короткой ссылки не существует")
	}

	return longURL, nil
}
