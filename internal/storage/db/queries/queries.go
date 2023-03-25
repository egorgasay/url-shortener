package queries

import (
	"database/sql"
	"errors"
)

type Query string
type Name int

const (
	InsertURL = iota
	GetLongLink
	FindMaxURL
	GetAllLinksByCookie
	MarkAsDeleted
	GetShortLink
)

var queriesSqlite3 = map[Name]Query{
	InsertURL:           "INSERT INTO urls (long, short, cookie) VALUES (?, ?, ?)",
	GetLongLink:         "SELECT long FROM urls WHERE short = ?",
	FindMaxURL:          "SELECT MAX(id) FROM urls",
	GetAllLinksByCookie: "SELECT short, long FROM urls WHERE cookie = ?",
	MarkAsDeleted:       "UPDATE urls SET deleted = 1 WHERE short = ? AND cookie = ?",
}

var queriesPostgres = map[Name]Query{
	InsertURL:           "INSERT INTO urls (long, short, cookie, deleted) VALUES ($1, $2, $3, false)",
	GetLongLink:         `SELECT long, deleted FROM urls WHERE short = $1`,
	FindMaxURL:          `SELECT MAX(id) FROM urls`,
	GetAllLinksByCookie: `SELECT short, long FROM urls WHERE cookie = $1`,
	MarkAsDeleted:       `UPDATE urls SET deleted = true WHERE short = $1 and cookie = $2`,
	GetShortLink:        "SELECT short FROM urls WHERE long = $1",
}

var queriesMySQL = map[Name]Query{
	InsertURL:           "INSERT INTO urls (`longURL`, `shortURL`, `cookie`) VALUES (?, ?, ?)",
	GetLongLink:         "SELECT `longURL` FROM urls WHERE `shortURL` = ?",
	FindMaxURL:          "SELECT MAX(id) FROM urls",
	GetAllLinksByCookie: "SELECT `short`, `long` FROM urls WHERE `cookie` = ?",
	MarkAsDeleted:       "UPDATE urls SET `deleted` = 1 WHERE `short` = ? AND `cookie` = ?",
}

var NotFoundError = errors.New("query not found")
var NilStatementError = errors.New("query statement is nil")

var statements = make(map[Name]*sql.Stmt, 10)

func Prepare(DB *sql.DB, vendor string) error {
	var queries map[Name]Query
	switch vendor {
	case "sqlite3":
		queries = queriesSqlite3
	case "postgres":
		queries = queriesPostgres
	case "mysql":
		queries = queriesMySQL
	}

	for n, q := range queries {
		prep, err := DB.Prepare(string(q))
		if err != nil {
			return err
		}
		statements[n] = prep
	}
	return nil
}

func GetPreparedStatement(name int) (*sql.Stmt, error) {
	stmt, ok := statements[Name(name)]
	if !ok {
		return nil, NotFoundError
	}

	if stmt == nil {
		return nil, NilStatementError
	}

	return stmt, nil
}
