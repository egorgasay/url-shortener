package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"io"
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

func (h Handler) GetLinkHandler(c *gin.Context) {
	longURL, err := h.services.GetLink.GetLink(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}
	c.Header("Location", longURL)
	c.Status(307)
}

func (h Handler) CreateLinkHandler(c *gin.Context) {
	b, err := io.ReadAll(c.Request.Body)
	if err != nil || len(b) < 3 {
		c.Error(errors.New("недопустимый URL"))
		c.AbortWithStatus(500)
		return
	}
	defer c.Request.Body.Close()
	shortURL, err := h.services.CreateLink.CreateLink(string(b))
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(400)
		return
	}
	c.Status(201)
	c.Writer.WriteString(domain + shortURL)
}
