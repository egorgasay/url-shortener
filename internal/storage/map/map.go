package mapstorage

import (
	"errors"
	"sync"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
)

type MapStorage struct {
	mu        sync.RWMutex
	container map[shortURL]data
}

func (s *MapStorage) MarkAsDeleted(shortURL, cookie string) {
	//TODO implement me
	panic("implement me")
}

const MapStorageType storage.Type = "map"

type shortURL string
type data struct {
	cookie  string
	longURL string
}

func NewMapStorage() storage.IStorage {
	db := make(map[shortURL]data, 10)
	return &MapStorage{container: db}
}

func (s *MapStorage) AddLink(longURL, ShortURL, cookie string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
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
