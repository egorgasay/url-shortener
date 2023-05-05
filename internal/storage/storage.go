package storage

import (
	"errors"
	"url-shortener/internal/schema"
)

// IStorage interface for a storage.
type IStorage interface {
	FindMaxID() (int, error)
	AddLink(longURL, shortURL, cookie string) (string, error)
	GetLongLink(shortURL string) (longURL string, err error)
	GetAllLinksByCookie(cookie string, baseURL string) (URLs []schema.URL, err error)
	Ping() error
	MarkAsDeleted(shortURL, cookie string) error
	Shutdown() error
	URLsCount() (int, error)
	UsersCount() (int, error)
}

// Type storage type.
type Type string

// ErrDeleted when URL was marked as deleted.
var ErrDeleted = errors.New("URL was marked as deleted")
