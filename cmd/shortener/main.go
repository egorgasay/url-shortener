package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"url-shortener/config"
	handlers "url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/routes"
)

func main() {
	cfg := config.New()

	storage, err := repository.New(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}

	router := gin.Default()
	handler := handlers.NewHandler(storage)
	public := router.Group("/")
	routes.PublicRoutes(public, handler)

	router.Use(gzip.Gzip(gzip.BestSpeed))
	router.Run(cfg.Host)
}
