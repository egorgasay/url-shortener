package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

const (
	alphabet    string = "AB1CDEFG2HIJKLM3NOPQRS4TUVW5XYZabc6defgh7ijklmn8opqrs9tuvw0xyz"
	lenAlphabet int    = 62
)

type CreateLinkSqlite struct {
	db *sql.DB
}

func NewCreateLinkSqlite(db *sql.DB) *CreateLinkSqlite {
	return &CreateLinkSqlite{db: db}
}

func (cr CreateLinkSqlite) CreateLink(longURL string) (string, error) {
	var li int

	stm := cr.db.QueryRow("SELECT MAX(id) FROM urls")

	err := stm.Scan(&li)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	shortURL := getShortName(li + 1)

	valueStrings := fmt.Sprintf("'%s','%s'", longURL, shortURL)
	stmt := fmt.Sprintf("INSERT INTO urls (long, short) VALUES (%s)", valueStrings) //nolint:gosec

	_, err = cr.db.Exec(stmt, valueStrings)

	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func getShortName(lastID int) (shrtURL string) {
	allNums := []int{}

	if lastID < 100000 {
		lastID = 10000 * lastID
	}

	for lastID > 0 {
		allNums = append(allNums, lastID%lenAlphabet)
		lastID /= lenAlphabet
	}

	// разворачиваем слайс
	for i, j := 0, len(allNums)-1; i < j; i, j = i+1, j-1 {
		allNums[i], allNums[j] = allNums[j], allNums[i]
	}

	chars := []string{}
	for _, el := range allNums {
		chars = append(chars, string(alphabet[el]))
	}

	shrtURL = strings.Join(chars, "")

	return
}
