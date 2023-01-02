package config

import (
	"os"
	"url-shortener/internal/repository"
)

var Domain = "http://127.0.0.1:8080/"

type Config struct {
	DBConfig *repository.Config
}

func New(baseURL string) *Config {
	if addr, ok := os.LookupEnv("BASE_URL"); ok {
		Domain = addr
	} else {
		Domain = baseURL
	}

	return &Config{
		DBConfig: &repository.Config{
			DriverName:     "file", // можно выбрать между map, sqlite3 и file
			DataSourceName: "urlshortener.txt",
		},
	}
}
