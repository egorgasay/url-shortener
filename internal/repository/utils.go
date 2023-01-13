package repository

import (
	"bufio"
	"database/sql"
	"errors"
	"os"
	"strings"
)

func IsDatabaseExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func InitDatabase(db *sql.DB) error {
	file, err := os.Open("schema.sql")
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
