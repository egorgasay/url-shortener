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

// AddLink adds a link to the repository.
func (s *MapStorage) AddLink(longURL, ShortURL, cookie string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.container[shortURL(ShortURL)]; ok {
		return ShortURL, service.ErrExists
	}

	s.container[shortURL(ShortURL)] = data{cookie: cookie, longURL: longURL}

	return ShortURL, nil
}

// FindMaxID gets len of the repository.
func (s *MapStorage) FindMaxID() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.container), nil
}

// GetLongLink gets a long link from the repository.
func (s *MapStorage) GetLongLink(ShortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, ok := s.container[shortURL(ShortURL)]
	if !ok {
		return "", errors.New("короткой ссылки не существует")
	}

	if record.deleted {
		return "", storage.ErrDeleted
	}

	return record.longURL, nil
}

// GetAllLinksByCookie gets all links ([]schema.URL) by cookie.
func (s *MapStorage) GetAllLinksByCookie(cookie, baseURL string) ([]schema.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var links []schema.URL

	for short, dt := range s.container {
		if dt.cookie == cookie {
			links = append(links, schema.URL{LongURL: dt.longURL, ShortURL: baseURL + string(short)})
		}
	}

	return links, nil
}

// MarkAsDeleted finds a URL and marks it as deleted.
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

// Ping checks connection with the repository.
func (s *MapStorage) Ping() error {
	if s.container == nil {
		return errors.New("хранилище не существует")
	}

	return nil
}
