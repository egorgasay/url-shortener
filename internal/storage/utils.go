package storage

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
)

func IsDBSqliteExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func InitDatabase(db *sql.DB, schema string) error {
	file, err := os.Open(schema)
	if err != nil {
		return err
	}
	defer file.Close()

	ctx := context.TODO()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var scanner = bufio.NewScanner(file)
	var queries strings.Builder

	for scanner.Scan() {
		queries.Write([]byte(scanner.Text()))
	}

	queriesArr := strings.Split(queries.String(), ";EOQ")

	for _, query := range queriesArr {
		_, err = tx.Exec(query)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
