package sqlite3

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"log"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage/db/service"

	_ "github.com/mattn/go-sqlite3"
)

const insertURL = "INSERT INTO urls (long, short, cookie) VALUES (?, ?, ?)"
const getLongLink = "SELECT long FROM urls WHERE short = ?"
const findMaxURL = "SELECT MAX(id) FROM urls"
const getAllLinksByCookie = "SELECT short, long FROM urls WHERE cookie = ?"
const markAsDeleted = "UPDATE urls SET deleted = 1 WHERE short = ? AND cookie = ?"

type Sqlite3 struct {
	DB *sql.DB
}

func New(db *sql.DB) service.IRealStorage {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/postgres",
		"sqlite", driver)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = m.Up()
	if err != nil {
		if err.Error() != "no change" {
			log.Fatal(err)
		}
	}

	return Sqlite3{DB: db}
}

func (s Sqlite3) AddLink(longURL, shortURL, cookie string) (string, error) {
	stmt, err := s.DB.Prepare(insertURL)
	if err != nil {
		return "", err
	}

	_, err = stmt.Exec(
		sql.Named("long", longURL).Value,
		sql.Named("short", shortURL).Value,
		sql.Named("cookie", cookie).Value,
	)

	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s Sqlite3) FindMaxID() (int, error) {
	var id int

	stmt, err := s.DB.Prepare(findMaxURL)
	if err != nil {
		return 0, nil
	}

	stm := stmt.QueryRow()
	err = stm.Scan(&id)

	return id, err
}

func (s Sqlite3) GetLongLink(shortURL string) (longURL string, err error) {
	stmt, err := s.DB.Prepare(getLongLink)
	if err != nil {
		return "", nil
	}

	stm := stmt.QueryRow(sql.Named("short", shortURL).Value)
	err = stm.Scan(&longURL)

	return longURL, err
}

func (s Sqlite3) MarkAsDeleted(shortURL, cookie string) {
	stmt, err := s.DB.Prepare(markAsDeleted)
	if err != nil {
		log.Println(err)
	}

	_, err = stmt.Exec(
		sql.Named("short", shortURL).Value,
		sql.Named("cookie", cookie).Value,
	)

	if err != nil {
		log.Println(err)
	}
}

func (s Sqlite3) GetAllLinksByCookie(cookie, baseURL string) ([]schema.URL, error) {
	stmt, err := s.DB.Prepare(getAllLinksByCookie)
	if err != nil {
		return nil, nil
	}

	stm, err := stmt.Query(sql.Named("cookie", cookie).Value)
	if err != nil {
		return nil, err
	}

	err = stm.Err()
	if err != nil {
		return nil, err
	}

	var URLs []schema.URL

	for stm.Next() {
		short, long := "", ""

		err = stm.Scan(&short, &long)
		if err != nil {
			return nil, err
		}

		URLs = append(URLs, schema.URL{LongURL: long, ShortURL: baseURL + short})
	}

	return URLs, err
}

func (s Sqlite3) Ping() error {
	ctx := context.TODO()
	return s.DB.PingContext(ctx)
}
