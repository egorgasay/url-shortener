package dbstorage

import "url-shortener/internal/storage"

type query string
type operations map[operation]query
type operation string

// Name of all operations
const insertURL operation = "insertURL"
const findMaxID operation = "findMaxID"
const getShortLink operation = "getShortLink"
const getLongLink operation = "getLongLink"

// Base Queries
const findMaxURLID query = "SELECT MAX(id) FROM urls"

// Queries for PostgreSQL
const insertURLPostgres query = "INSERT INTO urls (long, short, cookie) VALUES ($1, $2, $3)"
const getShortLinkPostgres query = "SELECT short FROM urls WHERE long = $1"
const getLongLinkPostgres query = "SELECT long FROM urls WHERE short = $1"

// Queries for MySQL
const insertURLMySQL query = "INSERT INTO urls (long, short, cookie) VALUES (?, ?, ?)"
const getLongLinkMySQL query = "SELECT long FROM urls WHERE short = ?"

// Queries for Sqlite3
const insertURLSqlite3 query = "INSERT INTO urls (long, short, cookie) VALUES (?, ?, ?)"
const getLongLinkSqlite3 query = "SELECT long FROM urls WHERE short = ?"

// queryStorage map that handle queries for all supported DBs
var queryStorage = map[storage.Type]operations{
	storage.Type("postgres"): operations{
		insertURL:    insertURLPostgres,
		findMaxID:    findMaxURLID,
		getShortLink: getShortLinkPostgres,
		getLongLink:  getLongLinkPostgres,
	},
	storage.Type("mysql"): operations{
		insertURL:   insertURLMySQL,
		findMaxID:   findMaxURLID,
		getLongLink: getLongLinkMySQL,
	},
	storage.Type("sqlite3"): operations{
		insertURL:   insertURLSqlite3,
		findMaxID:   findMaxURLID,
		getLongLink: getLongLinkSqlite3,
	},
}
