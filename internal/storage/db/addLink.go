package dbStorage

import (
	"github.com/speps/go-hashids"
)

const (
	alphabet string = "AB1CDEFG2HIJKLM3NOPQRS4TUVW5XYZabc6defgh7ijklmn8opqrs9tuvw0xyz"
)

func (s RealStorage) AddLink(longURL string, id int) (string, error) {
	shortURL := getShortName(id)
	stmt := "INSERT INTO urls (long, short) VALUES (?, ?)"

	_, err := s.Exec(stmt, longURL, shortURL)

	if err != nil {
		return "", err
	}

	return shortURL, err
}

func getShortName(lastID int) string {
	hd := hashids.NewData()
	hd.Salt = alphabet

	h, _ := hashids.NewWithData(hd)
	id, _ := h.Encode([]int{lastID})

	return id
}
