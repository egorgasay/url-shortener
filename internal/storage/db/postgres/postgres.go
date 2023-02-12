package postgres

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"log"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/db/service"
	shortenalgorithm "url-shortener/pkg/shortenAlgorithm"
)

const insertURL = "INSERT INTO urls (long, short, cookie, deleted) VALUES ($1, $2, $3, false)"
const getShortLink = "SELECT short FROM urls WHERE long = $1"
const getLongLink = `
SELECT long, deleted
FROM urls 
WHERE short = $1
`
const findMaxURL = "SELECT MAX(id) FROM urls"
const getAllLinksByCookie = "SELECT short, long FROM urls WHERE cookie = $1"
const markAsDeleted = "UPDATE urls SET deleted = true WHERE short = $1 AND cookie = $2"

type Postgres struct {
	DB *sql.DB
}

func New(db *sql.DB) service.IRealStorage {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/postgres",
		"postgres", driver)
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

	return Postgres{DB: db}
}

func (p Postgres) AddLink(longURL, shortURL, cookie string) (string, error) {
	stmt, err := p.DB.Prepare(insertURL)
	if err != nil {
		return "", err
	}

	_, err = stmt.Exec(
		sql.Named("long", longURL).Value,
		sql.Named("short", shortURL).Value,
		sql.Named("cookie", cookie).Value,
	)

	if err == nil {
		return shortURL, nil
	}

	e, ok := err.(*pq.Error)
	if !ok {
		log.Println("shouldn't be this ", err)
		return "", err
	}

	if e.Code != pgerrcode.UniqueViolation {
		return "", err
	}

	row := p.DB.QueryRow(getShortLink, sql.Named("long", longURL).Value)
	if row.Err() != nil {
		return "", err
	}

	err = row.Scan(&shortURL)
	if err == nil {
		return shortURL, service.ErrExists
	}

	log.Println("cycle")

	lastID, err := p.FindMaxID()
	if err != nil {
		return "", err
	}

	shortURL, err = shortenalgorithm.GetShortName(lastID + 1)
	if err != nil {
		return "", err
	}

	return p.AddLink(longURL, shortURL, cookie)
}

func (p Postgres) FindMaxID() (int, error) {
	var id int

	stmt, err := p.DB.Prepare(findMaxURL)
	if err != nil {
		return 0, nil
	}

	stm := stmt.QueryRow()
	err = stm.Scan(&id)

	return id, err
}

func (p Postgres) GetLongLink(shortURL string) (longURL string, err error) {
	stmt, err := p.DB.Prepare(getLongLink)
	if err != nil {
		return "", err
	}

	stm := stmt.QueryRow(sql.Named("short", shortURL).Value)
	var isDeleted bool
	err = stm.Scan(&longURL, &isDeleted)

	if isDeleted {
		return "", storage.ErrDeleted
	}

	return longURL, err
}

func (p Postgres) MarkAsDeleted(shortURL, cookie string) {
	stmt, err := p.DB.Prepare(markAsDeleted)
	if err != nil {
		log.Println(err)
	}

	_, err = stmt.Exec(shortURL, cookie)

	if err != nil {
		log.Println(err)
	}
}

func (p Postgres) GetAllLinksByCookie(cookie, baseURL string) ([]schema.URL, error) {
	stmt, err := p.DB.Prepare(getAllLinksByCookie)
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

func (p Postgres) Ping() error {
	ctx := context.TODO()
	return p.DB.PingContext(ctx)
}
