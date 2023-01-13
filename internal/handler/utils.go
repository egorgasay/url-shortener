package handler

import (
	"compress/gzip"
	"errors"
	"io"
	"net/url"
	"strings"
	"url-shortener/config"
)

func CreateLink(chars string) (*url.URL, error) {
	URL, err := url.Parse(*(config.F.BaseURL))
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
