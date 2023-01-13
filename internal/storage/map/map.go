package mapstorage

import (
	"sync"
	"url-shortener/internal/storage"
)

type MapStorage struct {
	mu        sync.RWMutex
	container map[string]string
}

func NewMapStorage() storage.IStorage {
	db := make(map[string]string, 10)
	return &MapStorage{container: db}
}
