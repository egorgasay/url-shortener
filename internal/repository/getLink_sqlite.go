package repository

import "url-shortener/internal/storage"

type GetLinkSqlite struct {
	db storage.IStorage
}

func NewGetLinkSqlite(db *storage.Storage) *GetLinkSqlite {
	if db == nil {
		panic("переменная storage равна nil")
	}

	return &GetLinkSqlite{db: db.DB}
}

func (gls GetLinkSqlite) GetLink(shortURL string) (longURL string, err error) {
	return gls.db.GetLongLink(shortURL)
}
