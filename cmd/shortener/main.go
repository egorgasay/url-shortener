package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"url-shortener/config"
	handlers "url-shortener/internal/handler"
	"url-shortener/internal/repository"
)

func main() {
	cfg := config.New()

	storage, err := repository.New(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}

	router := gin.Default()
	handler := handlers.NewHandler(storage)

	router.GET("/:id", handler.GetLinkHandler)
	router.POST("/", handler.CreateLinkHandler)
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
