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

// GetLink calls storage method GetLink.
// Returns a long URL and an error.
func (uc UseCase) GetLink(shortURL string) (longURL string, err error) {
	return uc.storage.GetLongLink(shortURL)
}

// MarkAsDeleted calls storage method MarkAsDeleted.
func (uc UseCase) MarkAsDeleted(shortURL, cookie string) {
	err := uc.storage.MarkAsDeleted(shortURL, cookie)
	if err != nil {
		// TODO: add zap logger
		log.Println("can't mark as deleted", err)
	}
}

// CreateLink calls FindMaxID, GetShortName and then calls AddLink storage method to save the link.
func (uc UseCase) CreateLink(longURL, cookie string, chars ...string) (string, error) {
	id, err := uc.storage.FindMaxID()
	if err != nil {
		log.Println("can't find max id", err)
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

// GetAllLinksByCookie calls storage method GetAllLinksByCookie and execute json from the response.
func (uc UseCase) GetAllLinksByCookie(cookie, baseURL string) (URLs string, err error) {
	links, err := uc.storage.GetAllLinksByCookie(cookie, baseURL)
	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(links, "", "    ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Ping checks connection with db.
func (uc UseCase) Ping() error {
	return uc.storage.Ping()
}

// Batch gets urls for process and processes every url in separate goroutines.
func (uc UseCase) Batch(batchURLs []schema.BatchURL, cookie, baseURL string) ([]schema.ResponseBatchURL, error) {
	var respJSON []schema.ResponseBatchURL
	var respJSONch = make(chan schema.ResponseBatchURL, len(batchURLs))
	var errorsCh = make(chan error)

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(200)
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
