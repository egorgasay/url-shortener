package mapstorage

import (
	"errors"
	"log"
	"sync"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/db/service"
)

// MapStorage struct with a map and mutex for concurent use.
type MapStorage struct {
	mu        sync.RWMutex
	container map[shortURL]data
}

// MapStorageType ...
const MapStorageType storage.Type = "map"

// shortURL ...
type shortURL string

// data ...
type data struct {
	cookie  string
	longURL string
	deleted bool
}

// NewMapStorage constructor for storage.IStorage with map implementation.
func NewMapStorage() storage.IStorage {
	db := make(map[shortURL]data, 10)
	return &MapStorage{container: db}
}

// AddLink ...
func (s *MapStorage) AddLink(longURL, ShortURL, cookie string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.container[shortURL(ShortURL)]; ok {
		return ShortURL, service.ErrExists
	}

	s.container[shortURL(ShortURL)] = data{cookie: cookie, longURL: longURL}

	return ShortURL, nil
}

// FindMaxID ...
func (s *MapStorage) FindMaxID() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.container), nil
}

// GetLongLink ...
func (s *MapStorage) GetLongLink(ShortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	Data, ok := s.container[shortURL(ShortURL)]
	if !ok {
		return "", errors.New("короткой ссылки не существует")
	}

	if Data.deleted {
		return "", storage.ErrDeleted
	}

	return Data.longURL, nil
}

// GetAllLinksByCookie ...
func (s *MapStorage) GetAllLinksByCookie(cookie, baseURL string) ([]schema.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var URLs []schema.URL

	for short, dt := range s.container {
		if dt.cookie == cookie {
			URLs = append(URLs, schema.URL{LongURL: dt.longURL, ShortURL: baseURL + string(short)})
		}
	}

	return URLs, nil
}

// MarkAsDeleted ...
func (s *MapStorage) MarkAsDeleted(ShortURL, cookie string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	Data, ok := s.container[shortURL(ShortURL)]
	if !ok {
		return
	}

	if cookie != Data.cookie {
		log.Println("wrong cookie")
		return
	}

	Data.deleted = true
	s.container[shortURL(ShortURL)] = Data
}

// Ping ...
func (s *MapStorage) Ping() error {
	if s.container == nil {
		return errors.New("хранилище не существует")
	}

	return nil
}
