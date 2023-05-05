package storage

import (
	"context"
	"errors"
	"url-shortener/internal/schema"
)

// IStorage interface for a storage.
type IStorage interface {
	FindMaxID(ctx context.Context) (int, error)
	AddLink(ctx context.Context, longURL, shortURL, cookie string) (string, error)
	GetLongLink(ctx context.Context, shortURL string) (longURL string, err error)
	GetAllLinksByCookie(ctx context.Context, cookie string, baseURL string) (URLs []schema.URL, err error)
	Ping(ctx context.Context) error
	MarkAsDeleted(shortURL, cookie string) error
	Shutdown() error
	URLsCount(ctx context.Context) (int, error)
	UsersCount(ctx context.Context) (int, error)
}

// Type storage type.
type Type string

// ErrDeleted when URL was marked as deleted.
var ErrDeleted = errors.New("URL was marked as deleted")
