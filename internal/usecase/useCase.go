package usecase

import (
	"encoding/json"
	"errors"
	"log"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage/db/service"
	shortenalgorithm "url-shortener/pkg/shortenAlgorithm"
)

func (uc UseCase) GetLink(shortURL string) (longURL string, err error) {
	return uc.storage.GetLongLink(shortURL)
}

func (uc UseCase) CreateLink(longURL, cookie string, chars ...string) (string, error) {
	id, err := uc.storage.FindMaxID()
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

	return uc.storage.AddLink(longURL, shortURL, cookie)
}

func (uc UseCase) GetAllLinksByCookie(shortURL, baseURL string) (URLs string, err error) {
	links, err := uc.storage.GetAllLinksByCookie(shortURL, baseURL)
	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(links, "", "    ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (uc UseCase) Ping() error {
	return uc.storage.Ping()
}

func (uc UseCase) Batch(batchURLs []schema.BatchURL, cookie, baseURL string) ([]schema.ResponseBatchURL, error) {
	var respJSON []schema.ResponseBatchURL
	for _, pair := range batchURLs {
		short, err := uc.CreateLink(pair.Original, cookie, pair.Chars)
		if err != nil && !errors.Is(err, service.ErrExists) {
			return nil, err
		}

		respJSON = append(respJSON, schema.ResponseBatchURL{Chars: pair.Chars,
			Shorted: baseURL + short})
	}

	return respJSON, nil
}
