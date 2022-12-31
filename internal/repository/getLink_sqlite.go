package repository

type GetLinkSqlite struct {
	db IStorage
}

func NewGetLinkSqlite(db *Storage) *GetLinkSqlite {
	if db == nil {
		panic("переменная storage равна nil")
	}

	return &GetLinkSqlite{db: db.DB}
}

func (gls GetLinkSqlite) GetLink(shortURL string) (longURL string, err error) {
	return gls.db.GetLongLink(shortURL)
}
