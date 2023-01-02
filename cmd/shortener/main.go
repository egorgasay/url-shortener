package main

import (
	"flag"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"url-shortener/config"
	handlers "url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/routes"
)

var (
	host    *string
	baseURL *string
	path    *string
)

func init() {
	host = flag.String("a", "localhost:8080", "-a=host")
	baseURL = flag.String("b", "http://localhost:8080/", "-b=URL")
	path = flag.String("f", "urlshortener.txt", "-f=path")
}

func main() {
	flag.Parse()

	cfg := config.New(*baseURL, *path)

	storage, err := repository.New(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}

	router := gin.Default()
	handler := handlers.NewHandler(storage)
	public := router.Group("/")
	routes.PublicRoutes(public, *handler)

	router.Use(gzip.Gzip(gzip.BestSpeed))

	serverAddress := *host
	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		serverAddress = addr
	}

	router.Run(serverAddress)
}
