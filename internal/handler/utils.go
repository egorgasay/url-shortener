package handler

import (
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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
