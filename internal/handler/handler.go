package handler

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
	"url-shortener/internal/service"
)

const (
	domain string = "http://127.0.0.1:8080/"
)

type Handler struct {
	services *service.Service
}

func NewHandler(storage *sql.DB) *Handler {
	return &Handler{service.NewService(storage)}
}

func (h Handler) GetLinkHandler(w http.ResponseWriter, r *http.Request) {
	longURL, err := h.services.GetLink.GetLink(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	w.Header().Add("Location", longURL)
	//fmt.Println(w.Header())
	w.WriteHeader(307)
}

func (h Handler) CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else if len(b) < 3 {
		err = errors.New("недопустимый URL")
		http.Error(w, err.Error(), 500)
	}
	defer r.Body.Close()
	shortURL, err := h.services.CreateLink.CreateLink(string(b))
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(201)
	w.Write([]byte(domain + shortURL))
}
