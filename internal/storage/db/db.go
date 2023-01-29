package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"

	"github.com/jackc/pgerrcode"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type RealStorage struct {
	DB     *sql.DB
	Vendor storage.Type
}

const DBStorageType storage.Type = "postgres"

func NewRealStorage(db *sql.DB, vendor storage.Type) storage.IStorage {
	return &RealStorage{DB: db, Vendor: vendor}
}

var ErrExists = errors.New("the shortened URL already exists")

func getQuery(vendor storage.Type, oper operation) (query, error) {
	v, ok := queryStorage[vendor]
	if !ok {
		return "", fmt.Errorf("operation do not implemented for %s", vendor)
	}

	q, ok := v[oper]
	if !ok {
		return "", fmt.Errorf("operation do not implemented for %s", vendor)
	}

	return q, nil
}

func (s RealStorage) AddLink(longURL, shortURL, cookie string) (string, error) {
	addLinkQuery, err := getQuery(s.Vendor, insertURL)
	if err != nil {
		return "", err
	}

	stmt, err := s.DB.Prepare(string(addLinkQuery))
	if err != nil {
		return "", err
	}

	_, err = stmt.Exec(
		sql.Named("long", longURL).Value,
		sql.Named("short", shortURL).Value,
		sql.Named("cookie", cookie).Value,
	)
	if err != nil {
		e, ok := err.(*pq.Error)
		if !ok {
			return "", err
		}

		if e.Code != pgerrcode.UniqueViolation {
			return "", err
		}

		getShortLinkQuery, err := getQuery(s.Vendor, getShortLink)
		if err != nil {
			return "", err
		}

		row := s.DB.QueryRow(string(getShortLinkQuery), sql.Named("long", longURL).Value)
		if row.Err() != nil {
			return "", err
		}

		err = row.Scan(&shortURL)
		if err != nil {
			return "", err
		}

		return shortURL, ErrExists
	}

	return shortURL, nil
}

func (s RealStorage) FindMaxID() (int, error) {
	var id int

	findMaxQuery, err := getQuery(s.Vendor, findMaxID)
	if err != nil {
		return 0, nil
	}

	stmt, err := s.DB.Prepare(string(findMaxQuery))
	if err != nil {
		return 0, nil
	}

	stm := stmt.QueryRow()
	err = stm.Scan(&id)

	return id, err
}

func (s RealStorage) GetLongLink(shortURL string) (longURL string, err error) {
	GetLongLinkQuery, err := getQuery(s.Vendor, getLongLink)
	if err != nil {
		return "", nil
	}

	stmt, err := s.DB.Prepare(string(GetLongLinkQuery))
	if err != nil {
		return "", nil
	}

	stm := stmt.QueryRow(sql.Named("short", shortURL).Value)
	err = stm.Scan(&longURL)

	return longURL, err
}

func (s RealStorage) GetAllLinksByCookie(cookie, baseURL string) ([]schema.URL, error) {
	GetAllLinksByCookieQuery, err := getQuery(s.Vendor, findMaxID)
	if err != nil {
		return nil, nil
	}

	stmt, err := s.DB.Prepare(string(GetAllLinksByCookieQuery))
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

func (s RealStorage) Ping() error {
	ctx := context.TODO()
	return s.DB.PingContext(ctx)
}
