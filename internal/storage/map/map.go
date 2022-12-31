package mapstorage

import "sync"

type MapStorage struct {
	mu        sync.RWMutex
	container map[string]string
}

func NewMapStorage() MapStorage {
	db := make(map[string]string, 10)
	return MapStorage{container: db}
}
