package usecase

import (
	"encoding/json"
	"log"
	"url-shortener/internal/storage"
)

func GetLink(repo storage.IStorage, shortURL string) (longURL string, err error) {
	return repo.GetLongLink(shortURL)
}

func CreateLink(repo storage.IStorage, longURL, cookie string) (string, error) {
	id, err := repo.FindMaxID()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return repo.AddLink(longURL, id+1, cookie)
}

func GetAllLinksByCookie(repo storage.IStorage, shortURL string) (URLs string, err error) {
	links, err := repo.GetAllLinksByCookie(shortURL)
	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(links, "", "    ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
