package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	handlers "url-shortener/internal/handler"
	"url-shortener/internal/repository"
	//"url-shortener/internal/storage"
)

func main() {
	//strg := storage.NewStorage()
	cfg := repository.Config{"sqlite3", "urlshortener.db"}
	strg, err := repository.NewSqliteDb(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}
	handler := handlers.NewHandler(strg)
	router := mux.NewRouter()
	router.HandleFunc("/{id}", handler.GetLinkHandler)
	router.HandleFunc("/", handler.CreateLinkHandler)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
