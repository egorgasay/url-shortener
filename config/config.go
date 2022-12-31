package config

import (
	"os"
	"url-shortener/internal/repository"
)

var Domain = "http://127.0.0.1:8080/"

type Config struct {
	DBConfig *repository.Config
}

func New() *Config {
	if addr, ok := os.LookupEnv("BASE_URL"); ok {
		Domain = addr
	}

	return &Config{
		DBConfig: &repository.Config{
			DriverName:     "sqlite3", // можно выбрать между map и sqlite3
			DataSourceName: "urlshortener.db",
		},
	}
}
