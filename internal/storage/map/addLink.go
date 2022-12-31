package mapstorage

import "github.com/speps/go-hashids"

const (
	alphabet string = "AB1CDEFG2HIJKLM3NOPQRS4TUVW5XYZabc6defgh7ijklmn8opqrs9tuvw0xyz"
)

func (s *MapStorage) AddLink(longURL string, id int) (string, error) {
	shortURL := getShortName(id)

	s.mu.RLock()
	defer s.mu.RUnlock()
	s.container[shortURL] = longURL

	return shortURL, nil
}

func getShortName(lastID int) string {
	hd := hashids.NewData()
	hd.Salt = alphabet

	h, _ := hashids.NewWithData(hd)
	id, _ := h.Encode([]int{lastID})

	return id
}
