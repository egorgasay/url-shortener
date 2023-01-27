package storage

import (
	"bufio"
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
)

func IsDBUsedBefore(driver, cred string) bool {
	db, err := sql.Open(driver, cred)
	if err != nil {
		log.Println(err)
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT short FROM urls")
	if err != nil {
		log.Println(err)
		return false
	}

	row := stmt.QueryRow()

	err = row.Err()
	if err != nil {
		log.Println(err)
		return false
	}

	var s string

	err = row.Scan(&s)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func IsDBSqliteExist(path string) bool {
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
