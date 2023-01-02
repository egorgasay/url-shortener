package handler

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"url-shortener/config"
	"url-shortener/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(storage service.IService) *Handler {
	if storage == nil {
		panic("переменная storage равна nil")
	}

	return &Handler{service.NewService(storage)}
}

func (h Handler) GetLinkHandler(c *gin.Context) {
	longURL, err := h.services.GetLink(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	GzipHandler(c)

	c.Header("Location", longURL)
	c.Status(http.StatusTemporaryRedirect)
}

func (h Handler) CreateLinkHandler(c *gin.Context) {
	reader, err := h.UseGzipHandler(c)
	if err != nil {
		return
	}

	b, err := io.ReadAll(reader)
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
	shortURL, err := h.services.CreateLink(link)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(config.Domain)
	if err != nil {
		return nil, err
	}

	u.Path = shortURL

	return u, nil
}

func (h Handler) UseGzipHandler(c *gin.Context) (io.Reader, error) {
	var reader io.Reader

	if str, ok := c.Get("Accept-Encoding"); ok && !strings.Contains(str.(string), "gzip") {
		gz, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)

			return nil, err
		}

		reader = gz
		defer gz.Close()
	} else {
		reader = c.Request.Body
	}

	return reader, nil
}

func (h Handler) APICreateLinkHandler(c *gin.Context) {
	reader, err := h.UseGzipHandler(c)
	if err != nil {
		return
	}

	b, err := io.ReadAll(reader)
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

func GzipHandler(c *gin.Context) {
	if str, ok := c.Get("Accept-Encoding"); ok && !strings.Contains(str.(string), "gzip") {
		return
	}

	c.Writer.Header().Set("Content-Encoding", "gzip")
	c.Writer.Header().Set("Content-Type", "application/x-gzip")
}
