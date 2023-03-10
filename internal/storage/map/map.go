package mapstorage

import (
	"errors"
	"log"
	"sync"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
)

type MapStorage struct {
	mu        sync.RWMutex
	container map[shortURL]data
}

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

const MapStorageType storage.Type = "map"

type shortURL string
type data struct {
	cookie  string
	longURL string
	deleted bool
}

func NewMapStorage() storage.IStorage {
	db := make(map[shortURL]data, 10)
	return &MapStorage{container: db}
}

func (s *MapStorage) AddLink(longURL, ShortURL, cookie string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.container[shortURL(ShortURL)] = data{cookie: cookie, longURL: longURL}

	return ShortURL, nil
}

func (s *MapStorage) FindMaxID() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.container), nil
}

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

func (s *MapStorage) Ping() error {
	if s.container == nil {
		return errors.New("хранилище не существует")
	}

	return nil
}
