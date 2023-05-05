package service

import (
	"context"
	"errors"
	"url-shortener/internal/schema"
)

// IRealStorage interface for the database storage.
type IRealStorage interface {
	AddLink(ctx context.Context, longURL, shortURL, cookie string) (string, error)
	FindMaxID(ctx context.Context) (int, error)
	GetLongLink(ctx context.Context, shortURL string) (longURL string, err error)
	GetAllLinksByCookie(ctx context.Context, cookie, baseURL string) ([]schema.URL, error)
	Ping(ctx context.Context) error
	MarkAsDeleted(shortURL, cookie string) error
	Shutdown() error
	URLsCount(ctx context.Context) (int, error)
	UsersCount(ctx context.Context) (int, error)
}

// ErrExists occurs when the shortened URL already exists.
var ErrExists = errors.New("the shortened URL already exists")
