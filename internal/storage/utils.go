package storage

import (
	"bufio"
	"database/sql"
	"errors"
	"os"
	"strings"
)

func IsDBUsedBefore(db *sql.DB) bool {
	stmt, err := db.Prepare("SELECT short FROM urls")
	if err != nil {
		return false
	}

	row := stmt.QueryRow()

	err = row.Err()
	if err != nil {
		return false
	}

	var s string

	err = row.Scan(&s)
	return err == nil
}

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

	var scanner = bufio.NewScanner(file)
	var queries strings.Builder

	for scanner.Scan() {
		queries.Write([]byte(scanner.Text()))
	}

	_, err = db.Exec(queries.String())
	if err != nil {
		return err
	}

	return nil
}
