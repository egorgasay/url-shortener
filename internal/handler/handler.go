package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"net/url"
	"url-shortener/config"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(storage *repository.Storage) *Handler {
	if storage == nil {
		panic("переменная storage равна nil")
	}

	return &Handler{service.NewService(storage)}
}

func (h Handler) GetLinkHandler(c *gin.Context) {
	longURL, err := h.services.GetLink.GetLink(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.Header("Location", longURL)
	c.Status(http.StatusTemporaryRedirect)
}

func (h Handler) CreateLinkHandler(c *gin.Context) {
	b, err := io.ReadAll(c.Request.Body)
	if err != nil || len(b) < 3 {
		c.Error(errors.New("недопустимый URL"))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	shortURL, err := h.services.CreateLink.CreateLink(string(b))
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	u, err := url.Parse(config.Domain + "chars")
	if err != nil {
		log.Fatal(err)
	}

	u.Path = shortURL

	c.Status(http.StatusCreated)

	c.Writer.WriteString(u.String())
}
