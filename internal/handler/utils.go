package handler

import (
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/url"
	"strings"
	"time"
)

func CreateLink(chars, baseURL string) (*url.URL, error) {
	URL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	URL.Path = chars

	return URL, nil
}

func UseGzip(body io.Reader, contentType string) (data []byte, err error) {
	if strings.Contains(contentType, "gzip") {
		data, err = DecompressGzip(body)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = io.ReadAll(body)
		if err != nil {
			return nil, err
		}
	}

	if len(string(data)) < 3 {
		return nil, errors.New("недопустимый URL")
	}

	return data, nil
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

func NewCookie(key []byte) string {
	h := hmac.New(sha256.New, key)
	src := []byte(fmt.Sprint(time.Now().UnixNano()))
	h.Write(src)

	return hex.EncodeToString(h.Sum(nil)) + "-" + hex.EncodeToString(src)
}

func getCookies(c *gin.Context) (cookie string, err error) {
	cookie = c.Request.Header.Get("Authorization")
	if cookie != "" {
		return cookie, nil
	}

	cookie, err = c.Cookie("token")
	if cookie != "" {
		return cookie, nil
	}

	cookie, err = c.Cookie("session")
	if cookie != "" {
		return cookie, nil
	}

	return cookie, errors.New("no cookies was provided")
}

func setCookies(c *gin.Context, host string, key []byte) (cookie string) {
	cookie = NewCookie(key)
	domain := strings.Split(host, ":")[0]
	c.SetCookie("token", cookie, 60*3600, "",
		domain, false, false)
	c.Header("Authorization", cookie)

	return cookie
}

func checkCookies(cookie string, key []byte) bool {
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
