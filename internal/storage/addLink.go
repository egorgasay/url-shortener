package storage

import "strings"

const (
	alphabet    string = "AB1CDEFG2HIJKLM3NOPQRS4TUVW5XYZabc6defgh7ijklmn8opqrs9tuvw0xyz"
	lenAlphabet int    = 62
)

func (s MapStorage) AddLink(longURL string, id int) (string, error) {
	shortURL := getShortName(id)
	s[shortURL] = longURL
	return shortURL, nil
}

func (s RealStorage) AddLink(longURL string, id int) (string, error) {
	shortURL := getShortName(id)
	stmt := "INSERT INTO urls (long, short) VALUES (?, ?)"

	_, err := s.Exec(stmt, longURL, shortURL)

	if err != nil {
		return "", err
	}

	return shortURL, err
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

	return shrtURL
}
