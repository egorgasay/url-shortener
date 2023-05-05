package mysql

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/db/basic"
	"url-shortener/internal/storage/db/queries"
	"url-shortener/internal/storage/db/service"

	_ "github.com/go-sql-driver/mysql"
)

var (
	_ storage.IStorage = (*MySQL)(nil)
)

// MySQL struct with *sql.DB instance.
// It has methods for working with URLs.
type MySQL struct {
	basic.DB
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

	return &MySQL{DB: basic.DB{DB: db}}
}

// FindMaxID gets len of the repository.
func (m *MySQL) FindMaxID(ctx context.Context) (int, error) {
	var id sql.NullInt32

	stmt, err := queries.GetPreparedStatement(queries.FindMaxURL)
	if err != nil {
		return 0, nil
	}

	stm := stmt.QueryRowContext(ctx)
	err = stm.Scan(&id)

	return int(id.Int32), err
}
