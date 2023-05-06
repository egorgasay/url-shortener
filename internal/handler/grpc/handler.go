package grpchandler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"net"
	"url-shortener/config"
	"url-shortener/internal/storage"
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

// Ping checks is alive db or not.
func (h *Handler) Ping(ctx context.Context, req *shortener.PingRequest) (*shortener.PingResponse, error) {
	return &shortener.PingResponse{}, nil
}

// Create creates shortened link.
func (h *Handler) Create(ctx context.Context, req *shortener.CreateRequest) (*shortener.CreateResponse, error) {
	charsForURL, err := h.logic.CreateLink(ctx, req.GetUrl(), "")
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

	return &shortener.CreateResponse{Shortened: URL.String()}, nil
}

// CreateApi creates shortened link and adds token(cookie).
func (h *Handler) CreateApi(ctx context.Context, req *shortener.CreateRequest) (*shortener.CreateResponse, error) {
	token, authenticated := getOrCreateToken(ctx, h.conf.Key)
	if !authenticated {
		setToken(ctx, token)
	}

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

	return &shortener.CreateResponse{Shortened: URL.String()}, nil
}

// Get gets original link.
func (h *Handler) Get(ctx context.Context, req *shortener.GetRequest) (*shortener.GetResponse, error) {
	URL, err := h.logic.GetLink(ctx, req.GetShortened())
	if err != nil {
		if errors.Is(err, storage.ErrDeleted) {
			return nil, status.Errorf(codes.Unavailable, "Link was deleted")
		}
		return nil, status.Errorf(codes.NotFound, "Link not found")
	}

	return &shortener.GetResponse{OriginalUrl: URL}, nil
}

// GetAll gets all original links by token.
func (h *Handler) GetAll(ctx context.Context, req *shortener.GetAllByCookieRequest) (*shortener.GetAllByCookieResponse, error) {
	token, authenticated := getOrCreateToken(ctx, h.conf.Key)
	if !authenticated {
		return &shortener.GetAllByCookieResponse{}, nil
	}

	URLs, err := h.logic.GetAllLinksByCookie(ctx, token, h.conf.BaseURL)
	if err != nil {
		return nil, err
	}

	return &shortener.GetAllByCookieResponse{Urls: URLs}, nil
}

// Delete deletes shortened link if token is valid and link was created by same user.
func (h *Handler) Delete(ctx context.Context, req *shortener.DeleteRequest) (*shortener.DeleteResponse, error) {
	token, _ := getOrCreateToken(ctx, h.conf.Key)

	urls := req.GetShortenedUrls()
	go func(token string, s []string) {
		for _, URL := range s {
			h.logic.MarkAsDeleted(URL, token)
		}
	}(token, urls)

	return &shortener.DeleteResponse{}, nil
}

// Batch creates shortened links and adds token(cookie).
func (h *Handler) Batch(ctx context.Context, req *shortener.BatchRequest) (*shortener.BatchResponse, error) {
	token, authenticated := getOrCreateToken(ctx, h.conf.Key)
	if !authenticated {
		setToken(ctx, token)
	}

	urls := req.GetUrls()
	resp, err := h.logic.Batch(ctx, urls, token, h.conf.BaseURL)
	if err != nil {
		return nil, err
	}

	return &shortener.BatchResponse{Urls: resp}, nil
}

// GetStats returns stats about urls and users.
func (h *Handler) GetStats(ctx context.Context, req *shortener.GetStatsRequest) (*shortener.GetStatsResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	} else if !h.conf.TrustedSubNetwork.Contains(net.ParseIP(p.Addr.String())) {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	data, err := h.logic.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Error while getting stats")
	}

	return &shortener.GetStatsResponse{Urls: int32(data.URLs), Users: int32(data.Users)}, nil
}
