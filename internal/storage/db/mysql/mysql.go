package mysql

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage/db/queries"
	"url-shortener/internal/storage/db/service"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL struct with *sql.DB instance.
// It has methods for working with URLs.
type MySQL struct {
	DB *sql.DB
}

// New MySQL struct constructor.
func New(db *sql.DB, path string) service.IRealStorage {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"mysql", driver)
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

	return MySQL{DB: db}
}

// AddLink adds a link to the repository.
func (m MySQL) AddLink(longURL, shortURL, cookie string) (string, error) {
	stmt, err := queries.GetPreparedStatement(queries.InsertURL)
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

// FindMaxID gets len of the repository.
func (m MySQL) FindMaxID() (int, error) {
	var id sql.NullInt32

	stmt, err := queries.GetPreparedStatement(queries.FindMaxURL)
	if err != nil {
		return 0, nil
	}

	stm := stmt.QueryRow()
	err = stm.Scan(&id)

	return int(id.Int32), err
}

// GetLongLink gets a long link from the repository.
func (m MySQL) GetLongLink(shortURL string) (longURL string, err error) {
	stmt, err := queries.GetPreparedStatement(queries.GetLongLink)
	if err != nil {
		return "", nil
	}

	stm := stmt.QueryRow(sql.Named("short", shortURL).Value)
	err = stm.Scan(&longURL)

	return longURL, err
}

// MarkAsDeleted finds a URL and marks it as deleted.
func (m MySQL) MarkAsDeleted(shortURL, cookie string) {
	stmt, err := queries.GetPreparedStatement(queries.MarkAsDeleted)
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

// GetAllLinksByCookie gets all links ([]schema.URL) by cookie.
func (m MySQL) GetAllLinksByCookie(cookie, baseURL string) ([]schema.URL, error) {
	stmt, err := queries.GetPreparedStatement(queries.GetAllLinksByCookie)
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

	var links []schema.URL

	for stm.Next() {
		short, long := "", ""

		err = stm.Scan(&short, &long)
		if err != nil {
			return nil, err
		}

		links = append(links, schema.URL{LongURL: long, ShortURL: baseURL + short})
	}

	return links, err
}

// Ping checks connection with the repository.
func (m MySQL) Ping() error {
	ctx := context.TODO()
	return m.DB.PingContext(ctx)
}
