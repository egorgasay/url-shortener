package repository

import (
	"database/sql"
)

//import handlers "url-shortener/internal/handler"

type GetLinkSqlite struct {
	db *sql.DB
}

func NewGetLinkSqlite(db *sql.DB) *GetLinkSqlite {
	return &GetLinkSqlite{db: db}
}

func (gls GetLinkSqlite) GetLink(shrt string) (longUrl string, err error) {
	stm := gls.db.QueryRow("SELECT long FROM urls WHERE short = ?", shrt[1:])
	err = stm.Scan(&longUrl)
	if err != nil {
		return
	}
	return
}
