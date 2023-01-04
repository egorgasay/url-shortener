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

	c.Header("Location", longURL)
	c.Status(http.StatusTemporaryRedirect)
}

func (h Handler) CreateLinkHandler(c *gin.Context) {

	//contentType := c.Request.Header.Get("Content-Type")
	//
	//bodyBytes, err := io.ReadAll(body)
	//if err != nil {
	//	c.Error(err)
	//	c.AbortWithStatus(http.StatusInternalServerError)
	//
	//	return
	//}

	data, err := h.UseGzipHandler(c.Request.Body, c.Request.Header.Get("Content-Type"))
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	u, err := h.CreateLink(string(data))
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

func (h Handler) UseGzipHandler(body io.Reader, contentType string) (data []byte, err error) {
	if strings.Contains(contentType, "gzip") {
		data, err = DecompressGzip(body)
		if err != nil {
			return nil, err
		}
	}

	if len(string(data)) < 3 {
		return nil, errors.New("недопустимый URL")
	}

	return data, nil
}

func (h Handler) APICreateLinkHandler(c *gin.Context) {
	b, err := h.UseGzipHandler(c.Request.Body, c.Request.Header.Get("Content-Type"))
	if err != nil {
		c.Error(err)
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

	c.Header("Content-Type", "application/json")
	c.Status(http.StatusCreated)

	c.Writer.Write(URL)
}

func DecompressGzip(body io.Reader) ([]byte, error) {
	gz, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}

	defer gz.Close()

	data, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//func GzipHandler(c *gin.Context) {
//	if str, ok := c.Get("Accept-Encoding"); ok && !strings.Contains(str.(string), "gzip") {
//		return
//	}
//
//	c.Writer.Header().Set("Content-Encoding", "gzip")
//	//c.Writer.Header().Set("Content-Type", "application/x-gzip")
//}
