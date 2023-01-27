package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"url-shortener/config"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
	"url-shortener/internal/usecase"
)

type Handler struct {
	storage storage.IStorage
	conf    *config.Config
}

func NewHandler(storage storage.IStorage, cfg *config.Config) *Handler {
	if storage == nil {
		panic("storage равен nil")
	}

	if cfg == nil {
		panic("конфиг равен nil")
	}

	return &Handler{storage: storage, conf: cfg}
}

func (h Handler) GetLinkHandler(c *gin.Context) {
	cookie, err := getCookies(c)
	if err != nil || !checkCookies(cookie, h.conf.Key) {
		log.Println("New cookie was created")
		setCookies(c, h.conf.Host, h.conf.Key)
	}

	longURL, err := usecase.GetLink(h.storage, c.Param("id"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.Header("Location", longURL)
	c.Status(http.StatusTemporaryRedirect)
}

func (h Handler) GetAllLinksHandler(c *gin.Context) {
	cookie, err := getCookies(c)
	if err != nil || !checkCookies(cookie, h.conf.Key) {
		cookie = setCookies(c, h.conf.Host, h.conf.Key)
	}

	URLs, err := usecase.GetAllLinksByCookie(h.storage, cookie, h.conf.BaseURL)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.Header("Content-Type", "application/json")

	if URLs == "null" {
		c.Status(http.StatusNoContent)
	} else {
		c.Status(http.StatusOK)
	}

	c.Writer.WriteString(URLs)
}

func (h Handler) CreateLinkHandler(c *gin.Context) {
	cookie, err := getCookies(c)
	if err != nil || !checkCookies(cookie, h.conf.Key) {
		cookie = setCookies(c, h.conf.Host, h.conf.Key)
	}

	data, err := UseGzip(c.Request.Body, c.Request.Header.Get("Content-Type"))
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	charsForURL, err := usecase.CreateLink(h.storage, string(data), cookie)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	URL, err := CreateLink(charsForURL, h.conf.BaseURL)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.Status(http.StatusCreated)

	c.Writer.WriteString(URL.String())
}

func (h Handler) APICreateLinkHandler(c *gin.Context) {
	cookie, err := getCookies(c)
	if err != nil || !checkCookies(cookie, h.conf.Key) {
		setCookies(c, h.conf.Host, h.conf.Key)
	}

	b, err := UseGzip(c.Request.Body, c.Request.Header.Get("Content-Type"))
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	var rj schema.RequestJSON

	err = json.Unmarshal(b, &rj)
	if err != nil {
		c.Error(errors.New("некорректный JSON"))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	charsForURL, err := usecase.CreateLink(h.storage, rj.URL, cookie)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	URL, err := CreateLink(charsForURL, h.conf.BaseURL)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	respJSON := schema.ResponseJSON{Result: URL.String()}

	rawURL, err := json.Marshal(respJSON)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.Header("Content-Type", "application/json")
	c.Status(http.StatusCreated)

	c.Writer.Write(rawURL)
}

func (h Handler) Ping(c *gin.Context) {
	err := usecase.Ping(h.storage)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
