package config

import (
	"flag"
	"os"
	"url-shortener/internal/repository"
)

const (
	mapStorage  string = "map"
	fileStorage string = "file"
	dbStorage   string = "sqlite3"
)

const (
	defaultURL     = "http://127.0.0.1:8080/"
	defaultHost    = "localhost:8080"
	defaultPath    = "urlshortener.txt"
	defaultStorage = mapStorage
)

type Flag struct {
	host    *string
	BaseURL *string
	path    *string
	storage *string
}

var F Flag

func init() {
	F.host = flag.String("a", defaultHost, "-a=host")
	F.BaseURL = flag.String("b", defaultURL, "-b=URL")
	F.path = flag.String("f", defaultPath, "-f=path")
	F.storage = flag.String("s", defaultStorage, "-s=storage")
}

type Config struct {
	Host     string
	DBConfig *repository.Config
}

func New() *Config {
	flag.Parse()

	if addr, ok := os.LookupEnv("BASE_URL"); ok {
		F.BaseURL = &addr
	}

	if fsp, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		F.path = &fsp
	}

	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		F.host = &addr
	}

	if *F.storage != mapStorage && *F.storage != fileStorage &&
		*F.storage != dbStorage {
		panic("Type of storage is not supported")
	}

	return &Config{
		Host: *F.host,
		DBConfig: &repository.Config{
			DriverName:     *F.storage,
			DataSourceName: *F.path,
		},
	}
}
