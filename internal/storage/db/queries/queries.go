package queries

import (
	"database/sql"
	"errors"
)

// Query text of query.
type Query string

// Name number of query.
type Name int

// Query names.
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
	FindMaxURL:          "SELECT MAX(`id`) FROM urls",
	GetAllLinksByCookie: "SELECT `shortURL`, `longURL` FROM urls WHERE `cookie` = ?",
	MarkAsDeleted:       "UPDATE urls SET `deleted` = 1 WHERE `shortURL` = ? AND `cookie` = ?",
}

// ErrNotFound occurs when query was not found.
var ErrNotFound = errors.New("the query was not found")

// ErrNilStatement occurs query statement is nil.
var ErrNilStatement = errors.New("query statement is nil")

var statements = make(map[Name]*sql.Stmt, 10)

// Prepare prepares all queries for db instance.
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

// GetPreparedStatement returns *sql.Stmt by name of query.
func GetPreparedStatement(name int) (*sql.Stmt, error) {
	stmt, ok := statements[Name(name)]
	if !ok {
		return nil, ErrNotFound
	}

	if stmt == nil {
		return nil, ErrNilStatement
	}

	return stmt, nil
}
