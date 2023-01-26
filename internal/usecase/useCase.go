package usecase

import (
	"encoding/json"
	"log"
	"url-shortener/internal/storage"
	shortenalgorithm "url-shortener/pkg/shortenAlgorithm"
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

	shortURL, err := shortenalgorithm.GetShortName(id + 1)
	if err != nil {
		return "", err
	}

	return repo.AddLink(longURL, shortURL, cookie)
}

func GetAllLinksByCookie(repo storage.IStorage, shortURL, baseURL string) (URLs string, err error) {
	links, err := repo.GetAllLinksByCookie(shortURL, baseURL)
	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(links, "", "    ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
