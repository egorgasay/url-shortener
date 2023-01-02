package repository

import (
	"log"
)

func (r Repository) CreateLink(longURL string) (string, error) {
	id, err := r.repo.FindMaxID()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return r.repo.AddLink(longURL, id+1)
}
