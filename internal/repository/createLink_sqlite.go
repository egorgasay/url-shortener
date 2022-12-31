package repository

import (
	"log"
	"url-shortener/internal/storage"
)

type CreateLinkSqlite struct {
	db storage.IStorage
}

func NewCreateLinkSqlite(db *storage.Storage) *CreateLinkSqlite {
	if db == nil {
		panic("переменная storage равна nil")
	}

	return &CreateLinkSqlite{db: db.DB}
}

func (cr CreateLinkSqlite) CreateLink(longURL string) (string, error) {
	id, err := cr.db.FindMaxID()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return cr.db.AddLink(longURL, id+1)
}
