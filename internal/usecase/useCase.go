package usecase

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage/db/service"
	shortener "url-shortener/pkg/api"
	shortenalgorithm "url-shortener/pkg/shortenAlgorithm"
)

// GetLink calls storage method GetLink.
// Returns a long URL and an error.
func (uc UseCase) GetLink(ctx context.Context, shortURL string) (longURL string, err error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	return uc.storage.GetLongLink(ctx, shortURL)
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
func (uc UseCase) CreateLink(ctx context.Context, longURL, cookie string, chars ...string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	id, err := uc.storage.FindMaxID(ctx)
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

	return uc.storage.AddLink(ctx, longURL, shortURL, cookie)
}

// GetAllLinksByCookie calls storage method GetAllLinksByCookie and execute json from the response.
func (uc UseCase) GetAllLinksByCookie(ctx context.Context, cookie, baseURL string) ([]*shortener.UserURL, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	links, err := uc.storage.GetAllLinksByCookie(ctx, cookie, baseURL)
	if err != nil {
		return nil, fmt.Errorf("can't get links by cookie: %w", err)
	}

	return links, nil
}

// Ping checks connection with db.
func (uc UseCase) Ping(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	return uc.storage.Ping(ctx)
}

// Batch gets urls for process and processes every url in separate goroutines.
func (uc UseCase) Batch(ctx context.Context, batchURLs []*shortener.LongAndShortURL, cookie, baseURL string) ([]*shortener.CharsAndShortURL, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var resp = make([]*shortener.CharsAndShortURL, len(batchURLs))

	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(200)

	for i, pair := range batchURLs {
		pair := pair
		i := i
		g.Go(func() error {
			short, err := uc.CreateLink(ctx, pair.OriginalUrl, cookie, pair.CorrelationId)
			if err != nil && !errors.Is(err, service.ErrExists) {
				return err
			}

			resp[i] = &shortener.CharsAndShortURL{CorrelationId: pair.CorrelationId, ShortUrl: baseURL + short}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("can't batch: %w", err)
	}

	return resp, nil
}

// GetStats calls storage method GetUser.
func (uc UseCase) GetStats(ctx context.Context) (stats schema.StatsResponse, err error) {
	if ctx.Err() != nil {
		return stats, ctx.Err()
	}

	stats.URLs, err = uc.storage.URLsCount(ctx)
	if err != nil {
		return stats, fmt.Errorf("can't get urls count: %w", err)
	}

	stats.Users, err = uc.storage.UsersCount(ctx)
	if err != nil {
		return stats, fmt.Errorf("can't get users count: %w", err)
	}

	return stats, nil
}
