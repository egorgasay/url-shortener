package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	handlers "url-shortener/internal/handler"
	"url-shortener/internal/repository"
)

func main() {
	cfg := repository.Config{DriverName: "sqlite3", DataSourceName: "urlshortener.db"}

	strg, err := repository.NewSqliteDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}

	router := gin.Default()
	handler := handlers.NewHandler(strg)

	router.GET("/:id", handler.GetLinkHandler)
	router.POST("/", handler.CreateLinkHandler)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", router))
}
