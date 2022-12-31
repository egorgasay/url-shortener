package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"net/url"
	"url-shortener/config"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"
)

type Handler struct {
	services *service.Service
}

func NewHandler(storage *storage.Storage) *Handler {
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

	u, err := h.CreateLink(string(b))
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.Status(http.StatusCreated)

	c.Writer.WriteString(u.String())
}

func (h Handler) CreateLink(link string) (*url.URL, error) {
	shortURL, err := h.services.CreateLink.CreateLink(link)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(config.Domain + "chars")
	if err != nil {
		return nil, err
	}

	u.Path = shortURL

	return u, nil
}

func (h Handler) APICreateLinkHandler(c *gin.Context) {
	b, err := io.ReadAll(c.Request.Body)
	if err != nil || len(b) < 3 {
		c.Error(errors.New("недопустимый URL"))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	type RequestJSON struct {
		URL string `json:"url"`
	}

	var rj RequestJSON

	err = json.Unmarshal(b, &rj)
	if err != nil {
		c.Error(errors.New("некорректный JSON"))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	u, err := h.CreateLink(rj.URL)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	type ResponseJSON struct {
		Result string `json:"result"`
	}

	respJSON := ResponseJSON{Result: u.String()}

	URL, err := json.Marshal(respJSON)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.Status(http.StatusCreated)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteString(string(URL))
}
