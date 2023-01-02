package config

import (
	"os"
	"url-shortener/internal/repository"
)

var Domain = "http://127.0.0.1:8080/"

type Config struct {
	DBConfig *repository.Config
}

func New(baseURL, path string) *Config {
	if addr, ok := os.LookupEnv("BASE_URL"); ok {
		Domain = addr
	} else {
		Domain = baseURL
	}

	if fsp, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		path = fsp
	}

	return &Config{
		DBConfig: &repository.Config{
			DriverName:     "file", // можно выбрать между map, sqlite3 и file
			DataSourceName: path,
		},
	}
}
