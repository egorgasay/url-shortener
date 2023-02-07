package config

import (
	"flag"
	"github.com/sethvargo/go-password/password"
	"log"
	"os"
	"url-shortener/internal/dockerdb"
	"url-shortener/internal/repository"
	"url-shortener/internal/storage"
	dbstorage "url-shortener/internal/storage/db"
	filestorage "url-shortener/internal/storage/file"
	mapstorage "url-shortener/internal/storage/map"
	getfreeport "url-shortener/pkg/getFreePort"
)

const (
	defaultURL     = "http://127.0.0.1:8080/"
	defaultHost    = "127.0.0.1:8080"
	defaultPath    = "urlshortener.txt"
	defaultStorage = dbstorage.DBStorageType
	defaultdsn     = ""
	defaultvdb     = ""
)

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

	if dsn, ok := os.LookupEnv("DATABASE_DSN"); ok {
		f.dsn = &dsn
	}

	if *f.dsn == "" && storage.Type(*f.storage) == defaultStorage && *f.vdb == "" {
		s := ""
		if *f.path != "" {
			s = string(filestorage.FileStorageType)
			f.storage = &s
		} else {
			s = string(mapstorage.MapStorageType)
			f.storage = &s
		}
	}

	generated, err := password.Generate(17, 5, 0, false, false)
	if err != nil {
		log.Fatal("Generate: ", err)
	}

	log.Println(*f.dsn, *f.path, *f.storage)
	var ddb *dockerdb.VDB

	if vdb := *f.vdb; vdb != "" {
		err := dockerdb.Pull(*f.storage)
		if err != nil {
			log.Fatal("Pull: ", err)
		}

		port, err := getfreeport.GetFreePort()
		if err != nil {
			log.Fatal("GetFreePort: ", err)
		}

		cfg := dockerdb.CustomDB{
			DB: dockerdb.DB{
				Name:     vdb,
				User:     "admin",
				Password: generated,
			},
			Port:   port,
			Vendor: *f.storage,
		}

		ddb, err = dockerdb.New(cfg)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &Config{
		Host:    *f.host,
		BaseURL: *f.baseURL,
		Key:     []byte(generated),
		DBConfig: &repository.Config{
			DriverName:     storage.Type(*f.storage),
			DataSourcePath: *f.path,
			DataSourceCred: *f.dsn,
			VDB:            ddb,
		},
	}
}
