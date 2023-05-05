package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
	"url-shortener/config"
	"url-shortener/internal/storage/db/service"
	"url-shortener/internal/usecase"
	shortener "url-shortener/pkg/api"
)

// Handler struct that contains link to the logic layer and conf.
// It has methods for processing requests.
type Handler struct {
	conf  *config.Config
	logic usecase.UseCase
	shortener.UnimplementedShortenerServer
}

// New returns new Handler.
func NewHandler(conf *config.Config, logic usecase.UseCase) *Handler {
	return &Handler{
		conf:  conf,
		logic: logic,
	}
}

func (h *Handler) Ping(ctx context.Context, req *shortener.PingRequest) (*shortener.PingResponse, error) {
	return &shortener.PingResponse{}, nil
}

func (h *Handler) Create(ctx context.Context, req *shortener.CreateRequest) (*shortener.CreateResponse, error) {
	token, unauthenticated := getOrCreateToken(ctx, h.conf.Key)

	charsForURL, err := h.logic.CreateLink(ctx, req.GetUrl(), token)
	if err != nil {
		if !errors.Is(err, service.ErrExists) {
			return nil, err
		}

		URL, err := CreateLink(charsForURL, h.conf.BaseURL)
		if err != nil {
			return nil, err
		}

		return &shortener.CreateResponse{Shortened: URL.String()}, status.Errorf(codes.AlreadyExists, "Link already exists")
	}

	URL, err := CreateLink(charsForURL, h.conf.BaseURL)
	if err != nil {
		return nil, err
	}

	if unauthenticated {
		return &shortener.CreateResponse{Shortened: URL.String()}, status.Errorf(codes.Unauthenticated, token)
	}

	return &shortener.CreateResponse{Shortened: URL.String()}, nil
}

// CreateLink accepts chars and baseURL for building url.URL.
func CreateLink(chars, baseURL string) (*url.URL, error) {
	URL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	URL.Path = chars

	return URL, nil
}
