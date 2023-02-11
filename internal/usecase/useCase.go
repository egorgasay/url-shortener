package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/sync/errgroup"
	"log"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage/db/service"
	shortenalgorithm "url-shortener/pkg/shortenAlgorithm"
)

func (uc UseCase) GetLink(shortURL string) (longURL string, err error) {
	return uc.storage.GetLongLink(shortURL)
}

func (uc UseCase) MarkAsDeleted(shortURL, cookie string) {
	uc.storage.MarkAsDeleted(shortURL, cookie)
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
	var respJSONch = make(chan schema.ResponseBatchURL, len(batchURLs))
	var errorsCh = make(chan error)

	g, _ := errgroup.WithContext(context.Background())
	go func() {
		errorsCh <- g.Wait()
	}()

	for _, pair := range batchURLs {
		pair := pair
		g.Go(func() error {
			short, err := uc.CreateLink(pair.Original, cookie, pair.Chars)
			if err != nil && !errors.Is(err, service.ErrExists) {
				return err
			}

			respJSONch <- schema.ResponseBatchURL{Chars: pair.Chars,
				Shorted: baseURL + short}

			return nil
		})
	}

	i := 1
	for resp := range respJSONch {
		select {
		case err := <-errorsCh:
			if err != nil {
				return nil, err
			}
		default:
		}

		respJSON = append(respJSON, resp)

		if i == len(batchURLs) {
			close(respJSONch)
		}

		i++
	}

	return respJSON, nil
}
