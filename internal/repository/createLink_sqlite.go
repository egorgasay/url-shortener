package repository

import (
	"log"
)

type CreateLinkSqlite struct {
	db IStorage
}

func NewCreateLinkSqlite(db *Storage) *CreateLinkSqlite {
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

	shortURL, err := cr.db.AddLink(longURL, id+1)

	return shortURL, nil
}
