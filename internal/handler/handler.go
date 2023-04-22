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
	"url-shortener/internal/storage/db/service"
	"url-shortener/internal/usecase"
)

// Handler struct that contains link to the logic layer and conf.
// It has methods for processing requests.
type Handler struct {
	conf  *config.Config
	logic usecase.UseCase
}

// NewHandler creates an instance of the Handler.
func NewHandler(cfg *config.Config, logic usecase.UseCase) *Handler {
	if cfg == nil {
		panic("конфиг равен nil")
	}

	return &Handler{conf: cfg, logic: logic}
}

// GetLinkHandler accepts short url through the characters in the url (after the slash),
// returns a redirect to the URL that was shortened.
func (h Handler) GetLinkHandler(c *gin.Context) {
	longURL, err := h.logic.GetLink(c.Param("id"))
	if err != nil {
		log.Println(err)
		if errors.Is(err, storage.ErrDeleted) {
			c.AbortWithStatus(http.StatusGone)
			return
		}

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Header("Location", longURL)
	c.Status(http.StatusTemporaryRedirect)
}

// GetAllLinksHandler returns all URLs that have been shortened by a specific user,
// which is determined using a cookie provided upon request.
func (h Handler) GetAllLinksHandler(c *gin.Context) {
	cookie, err := getCookies(c)
	if err != nil || !checkCookies(cookie, h.conf.Key) {
		cookie = setCookies(c, h.conf.Key)
	}

	URLs, err := h.logic.GetAllLinksByCookie(cookie, h.conf.BaseURL)
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

// CreateLinkHandler accepts original link in the request (as plain text) and
// returns a shortened equivalent.
func (h Handler) CreateLinkHandler(c *gin.Context) {
	cookie, err := getCookies(c)
	if err != nil || !checkCookies(cookie, h.conf.Key) {
		cookie = setCookies(c, h.conf.Key)
	}

	data, err := UseGzip(c.Request.Body, c.Request.Header.Get("Content-Type"))
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	charsForURL, err := h.logic.CreateLink(string(data), cookie)
	if err != nil {
		if !errors.Is(err, service.ErrExists) {
			c.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusConflict)

		URL, err := CreateLink(charsForURL, h.conf.BaseURL)
		if err != nil {
			c.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		c.Writer.WriteString(URL.String())
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

// APICreateLinkHandler accepts original link in the request (as json) and
// returns a shortened equivalent.
func (h Handler) APICreateLinkHandler(c *gin.Context) {
	cookie, err := getCookies(c)
	if err != nil || !checkCookies(cookie, h.conf.Key) {
		cookie = setCookies(c, h.conf.Key)
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

	var isConflict bool
	charsForURL, err := h.logic.CreateLink(rj.URL, cookie)
	if err != nil {
		if !errors.Is(err, service.ErrExists) {
			c.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		isConflict = true
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

	if isConflict {
		c.Status(http.StatusConflict)
	} else {
		c.Status(http.StatusCreated)
	}

	c.Header("Content-Type", "application/json")
	c.Writer.Write(rawURL)
}

// Ping checks the connection to the database.
func (h Handler) Ping(c *gin.Context) {
	err := h.logic.Ping()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// BatchHandler accepts a batch of URLs and saves them.
// Returns correlation id and shortened urls in the response.
func (h Handler) BatchHandler(c *gin.Context) {
	cookie, err := getCookies(c)
	if err != nil || !checkCookies(cookie, h.conf.Key) {
		cookie = setCookies(c, h.conf.Key)
	}

	var batchURLs []schema.BatchURL
	err = c.BindJSON(&batchURLs)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	data, err := h.logic.Batch(batchURLs, cookie, h.conf.BaseURL)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Content-Type", "application/json")
	c.IndentedJSON(http.StatusCreated, data)
}

// APIDeleteLinksHandler accepts a batch of URLs and marks them as deleted.
func (h Handler) APIDeleteLinksHandler(c *gin.Context) {
	cookie, err := getCookies(c)
	if err != nil || !checkCookies(cookie, h.conf.Key) {
		cookie = setCookies(c, h.conf.Key)
	}

	var s []string
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}

	go func(cookie string, s []string) {
		for _, URL := range s {
			h.logic.MarkAsDeleted(URL, cookie)
		}
	}(cookie, s)

	c.Status(http.StatusAccepted)
	c.Header("Content-Type", "application/json")
}
