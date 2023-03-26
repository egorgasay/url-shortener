package config

import (
	"context"
	"flag"
	"github.com/egorgasay/dockerdb"
	"log"
	"os"
	"url-shortener/internal/repository"
	"url-shortener/internal/storage"
	dbstorage "url-shortener/internal/storage/db"
	filestorage "url-shortener/internal/storage/file"
	mapstorage "url-shortener/internal/storage/map"
)

const (
	defaultURL     = "http://127.0.0.1:8080/"
	defaultHost    = "127.0.0.1:8080"
	defaultPath    = "urlshortener.txt"
	defaultStorage = dbstorage.DBStorageType
	defaultdsn     = ""
	defaultvdb     = ""
)

// Flag struct for parsing from env and cmd args.
type Flag struct {
	host    *string
	baseURL *string
	path    *string
	storage *string
	dsn     *string
	vendor  *string
	vdb     *string
}

var f Flag

func init() {
	f.host = flag.String("a", defaultHost, "-a=host")
	f.baseURL = flag.String("b", defaultURL, "-b=URL")
	f.path = flag.String("f", defaultPath, "-f=path")
	f.storage = flag.String("s", string(defaultStorage), "-s=storage")
	f.dsn = flag.String("d", defaultdsn, "-d=connection_string")
	f.vdb = flag.String("vdb", defaultvdb, "-vdb=virtual_db_name")
}

// Config contains all the settings for configuring the application.
type Config struct {
	Host     string
	BaseURL  string
	Key      []byte
	DBConfig *repository.Config
}

// New initializing the config for the application.
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

	if dsn, ok := os.LookupEnv("DATABASE_DSN"); ok {
		f.dsn = &dsn
	}

	if *f.dsn == "" && storage.Type(*f.storage) == defaultStorage && *f.vdb == "" {
		s := ""
		if *f.path != defaultPath {
			s = string(filestorage.FileStorageType)
			f.storage = &s
		} else {
			s = string(mapstorage.MapStorageType)
			f.storage = &s
		}
	}

	//log.Println(*f.dsn, *f.path, *f.storage)
	var ddb *dockerdb.VDB
	var vdb = *f.vdb

	if vdb != "" {
		ctx := context.TODO()

		cfg := dockerdb.CustomDB{
			DB: dockerdb.DB{
				Name:     vdb,
				User:     "admin",
				Password: "admin",
			},
			Port: "12522",
			Vendor: dockerdb.Vendor{
				Name:  *f.storage,
				Image: *f.storage,
			},
		}

		var err error
		ddb, err = dockerdb.New(ctx, cfg)
		if err != nil {
			log.Fatal(err)
		}
		f.dsn = &ddb.ConnString
	}

	return &Config{
		Host:    *f.host,
		BaseURL: *f.baseURL,
		Key:     []byte("CHANGE ME"),
		DBConfig: &repository.Config{
			DriverName:     storage.Type(*f.storage),
			DataSourcePath: *f.path,
			DataSourceCred: *f.dsn,
			VDB:            ddb,
			Name:           vdb,
		},
	}
}
