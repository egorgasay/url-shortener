package service

import (
	"url-shortener/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type GetLink interface {
	GetLink(string) (string, error)
}

type CreateLink interface {
	CreateLink(string) (string, error)
}

type Service struct {
	GetLink
	CreateLink
}

func NewService(db *repository.Storage) *Service {
	if db == nil {
		panic("переменная storage равна nil")
	}

	return &Service{
		repository.NewGetLinkSqlite(db),
		repository.NewCreateLinkSqlite(db),
	}
}
