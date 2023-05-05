package mapstorage

import (
	"context"
	"errors"
	"sync"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/db/service"
)

var (
	_ storage.IStorage = (*MapStorage)(nil)
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
func (s *MapStorage) AddLink(ctx context.Context, longURL, ShortURL, cookie string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.container[shortURL(ShortURL)]; ok {
		return ShortURL, service.ErrExists
	}

	s.container[shortURL(ShortURL)] = data{cookie: cookie, longURL: longURL}

	return ShortURL, nil
}

// FindMaxID gets len of the repository.
func (s *MapStorage) FindMaxID(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.container), nil
}

// GetLongLink gets a long link from the repository.
func (s *MapStorage) GetLongLink(ctx context.Context, ShortURL string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

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
func (s *MapStorage) GetAllLinksByCookie(ctx context.Context, cookie, baseURL string) ([]schema.URL, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

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
func (s *MapStorage) MarkAsDeleted(ShortURL, cookie string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	Data, ok := s.container[shortURL(ShortURL)]
	if !ok {
		return errors.New("not found")
	}

	if cookie != Data.cookie {
		return errors.New("wrong cookie")
	}

	Data.deleted = true
	s.container[shortURL(ShortURL)] = Data

	return nil
}

// Ping checks connection with the repository.
func (s *MapStorage) Ping(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if s.container == nil {
		return errors.New("хранилище не существует")
	}

	return nil
}

// Shutdown clears the repository.
func (s *MapStorage) Shutdown() error {
	s.container = make(map[shortURL]data)
	return nil
}

// URLsCount gets count of the repository.
func (s *MapStorage) URLsCount(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.container), nil
}

// UsersCount gets count of users.
func (s *MapStorage) UsersCount(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var users = make(map[string]struct{}, 100)

	for _, dt := range s.container {
		users[dt.cookie] = struct{}{}
	}

	return len(users), nil
}
