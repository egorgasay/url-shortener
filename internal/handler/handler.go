package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"url-shortener/config"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
	"url-shortener/internal/usecase"

	_ "github.com/mattn/go-sqlite3"
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
	cookie, err := GetCookies(c)
	if err != nil || !CheckCookies(cookie, h.conf.Key) {
		log.Println("New cookie was created")
		SetCookies(c, h.conf.Host, h.conf.Key)
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
	cookie, err := GetCookies(c)
	if err != nil || !CheckCookies(cookie, h.conf.Key) {
		cookie = SetCookies(c, h.conf.Host, h.conf.Key)
	}

	URLs, err := usecase.GetAllLinksByCookie(h.storage, cookie)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}
	if URLs == "null" {
		c.Status(http.StatusNoContent)
	} else {
		c.Status(http.StatusOK)
	}

	c.Writer.WriteString(URLs)
}

func (h Handler) CreateLinkHandler(c *gin.Context) {
	cookie, err := GetCookies(c)
	if err != nil || !CheckCookies(cookie, h.conf.Key) {
		cookie = SetCookies(c, h.conf.Host, h.conf.Key)
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
	cookie, err := GetCookies(c)
	if err != nil || !CheckCookies(cookie, h.conf.Key) {
		SetCookies(c, h.conf.Host, h.conf.Key)
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

	URL, err := CreateLink(charsForURL, h.conf.Host)
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

func GetCookies(c *gin.Context) (string, error) {
	cookie, err := c.Cookie("token")
	if err != nil {
		return "", err
	}

	return cookie, err
}

func SetCookies(c *gin.Context, host string, key []byte) (cookie string) {
	cookie = NewCookie(key)
	domain := strings.Split(host, ":")[0]
	c.SetCookie("token", cookie, 10, "",
		domain, true, false)

	return cookie
}

func CheckCookies(cookie string, key []byte) bool {
	arr := strings.Split(cookie, "-")
	k, v := arr[0], arr[1]

	sign, err := hex.DecodeString(k)
	if err != nil {
		return false
	}

	data, err := hex.DecodeString(v)
	if err != nil {
		return false
	}

	h := hmac.New(sha256.New, key)
	h.Write(data)

	return hmac.Equal(sign, h.Sum(nil))
}
