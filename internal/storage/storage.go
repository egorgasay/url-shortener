package storage

import "url-shortener/internal/schema"

type IStorage interface {
	FindMaxID() (int, error)
	AddLink(longURL, shortURL, cookie string) (string, error)
	GetLongLink(shortURL string) (longURL string, err error)
	GetAllLinksByCookie(cookie string, baseURL string) (URLs []schema.URL, err error)
}

type Type string
