package config

import (
	"flag"
	"os"
	"url-shortener/internal/repository"
	"url-shortener/internal/storage"
	dbstorage "url-shortener/internal/storage/db"
	fileStorage "url-shortener/internal/storage/file"
	mapStorage "url-shortener/internal/storage/map"
)

const (
	defaultURL     = "http://127.0.0.1:8080/"
	defaultHost    = "127.0.0.1:8080"
	defaultPath    = "urlshortener.txt"
	defaultStorage = mapStorage.MapStorageType
)

type Flag struct {
	host    *string
	baseURL *string
	path    *string
	storage storage.Type
}

var f Flag

func init() {
	f.host = flag.String("a", defaultHost, "-a=host")
	f.baseURL = flag.String("b", defaultURL, "-b=URL")
	f.path = flag.String("f", defaultPath, "-f=path")
	f.storage = storage.Type(*flag.String("s", string(defaultStorage), "-s=storage"))
}

type Config struct {
	Host     string
	BaseURL  string
	Key      []byte
	DBConfig *repository.Config
}

func New() *Config {
	flag.Parse()

	if addr, ok := os.LookupEnv("BASE_URL"); ok {
		f.baseURL = &addr
	}

	if fsp, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		f.path = &fsp
	}

	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		f.host = &addr
	}

	if f.storage != mapStorage.MapStorageType && f.storage != fileStorage.FileStorageType &&
		f.storage != dbstorage.DBStorageType {
		panic("Type of storage is not supported")
	}

	return &Config{
		Host:    *f.host,
		BaseURL: *f.baseURL,
		Key:     []byte("CHANGE ME"),
		DBConfig: &repository.Config{
			DriverName:     f.storage,
			DataSourceName: *f.path,
		},
	}
}
