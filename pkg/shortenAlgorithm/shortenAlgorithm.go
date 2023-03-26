package shortenalgorithm

import "github.com/speps/go-hashids"

const (
	alphabet string = "AB1CDEFG2HIJKLM3NOPQRS4TUVW5XYZabc6defgh7ijklmn8opqrs9tuvw0xyz"
)

// GetShortName generates a short string equivalent for digit.
func GetShortName(lastID int) (string, error) {
	hd := hashids.NewData()
	hd.Salt = alphabet

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}

	id, err := h.Encode([]int{lastID})
	if err != nil {
		return "", err
	}

	return id, nil
}
