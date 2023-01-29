package usecase

import (
	"encoding/json"
	"errors"
	"log"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
	dbstorage "url-shortener/internal/storage/db"
	shortenalgorithm "url-shortener/pkg/shortenAlgorithm"
)

func GetLink(repo storage.IStorage, shortURL string) (longURL string, err error) {
	return repo.GetLongLink(shortURL)
}

func CreateLink(repo storage.IStorage, longURL, cookie string, chars ...string) (string, error) {
	id, err := repo.FindMaxID()
	if err != nil {
		log.Println(err)
		return "", err
	}

	var shortURL string
	if len(chars) > 0 {
		shortURL = chars[0]
	} else {
		shortURL, err = shortenalgorithm.GetShortName(id + 1)
		if err != nil {
			return "", err
		}
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

func Ping(repo storage.IStorage) error {
	return repo.Ping()
}

func Batch(repo storage.IStorage, batchURLs []schema.BatchURL, cookie, baseURL string) ([]schema.ResponseBatchURL, error) {
	var respJSON []schema.ResponseBatchURL
	for _, pair := range batchURLs {
		short, err := CreateLink(repo, pair.Original, cookie, pair.Chars)
		if err != nil && !errors.Is(err, dbstorage.ErrExists) {
			return nil, err
		}

		respJSON = append(respJSON, schema.ResponseBatchURL{Chars: pair.Chars,
			Shorted: baseURL + short})
	}

	return respJSON, nil
}
