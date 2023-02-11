package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
	"url-shortener/config"
	handlers "url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/routes"
	"url-shortener/internal/usecase"
)

func main() {
	cfg := config.New()

	storage, err := repository.New(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}

	logic := usecase.New(storage)
	router := gin.Default()
	handler := handlers.NewHandler(cfg, logic)

	public := router.Group("/")
	routes.PublicRoutes(public, handler)

	router.Use(gzip.Gzip(gzip.BestSpeed))

	go func() {
		router.Run(cfg.Host)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown Server ...")
}
