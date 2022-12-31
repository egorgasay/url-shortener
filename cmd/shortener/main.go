package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
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
	routes.PublicRoutes(public, *handler)

	serverAddress := "127.0.0.1"
	if addr, ok := os.LookupEnv("BASE_URL"); ok {
		serverAddress = addr
	}

	router.Run(serverAddress + ":8080")
}
