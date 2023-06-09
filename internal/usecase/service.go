package usecase

import (
	"url-shortener/internal/storage"
)

// UseCase logic layer main struct.
type UseCase struct {
	storage storage.IStorage
}

// New the UseCase struct builder.
func New(storage storage.IStorage) UseCase {
	return UseCase{storage: storage}
}
