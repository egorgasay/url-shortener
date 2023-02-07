package usecase

import "url-shortener/internal/storage"

type UseCase struct {
	storage storage.IStorage
}

func New(storage storage.IStorage) UseCase {
	return UseCase{storage: storage}
}
