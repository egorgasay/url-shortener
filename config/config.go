package config

import "url-shortener/internal/repository"

const (
	Domain string = "http://127.0.0.1:8080/"
)

type Config struct {
	DBConfig *repository.Config
}

func New() *Config {
	return &Config{
		DBConfig: &repository.Config{
			DriverName:     "sqlite3",
			DataSourceName: "urlshortener.db",
		},
	}
}
