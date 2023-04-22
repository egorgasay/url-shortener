package service

import (
	"errors"
	"url-shortener/internal/schema"
)

// IRealStorage interface for the database storage.
type IRealStorage interface {
	AddLink(longURL, shortURL, cookie string) (string, error)
	FindMaxID() (int, error)
	GetLongLink(shortURL string) (longURL string, err error)
	GetAllLinksByCookie(cookie, baseURL string) ([]schema.URL, error)
	Ping() error
	MarkAsDeleted(shortURL, cookie string) error
	Shutdown() error
}

// ErrExists occurs when the shortened URL already exists.
var ErrExists = errors.New("the shortened URL already exists")
