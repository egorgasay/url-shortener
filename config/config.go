package config

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/egorgasay/dockerdb"
	"io"
	"log"
	"net"
	"os"
	"reflect"
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
)

// Flag struct for parsing from env and cmd args.
type Flag struct {
	Host              *string `json:"server_address,omitempty"`
	BaseURL           *string `json:"base_url,omitempty"`
	Path              *string `json:"file_storage_path,omitempty"`
	Storage           *string `json:"storage,omitempty"`
	DSN               *string `json:"database_dsn,omitempty"`
	VDB               *string `json:"vdb,omitempty"`
	Cfg               *string
	Config            *string
	HTTPS             *bool   `json:"enable_https,omitempty"`
	TrustedSubNetwork *string `json:"trusted_subnet"`
}

var f Flag

// defaults for properly working the reflection.
var defaults = map[string]string{
	"Host":    defaultHost,
	"BaseURL": defaultURL,
	"Path":    defaultPath,
	"Storage": string(defaultStorage),
}

func init() {
	f.Host = flag.String("a", defaults["Host"], "-a=host")
	f.BaseURL = flag.String("b", defaults["BaseURL"], "-b=URL")
	f.Path = flag.String("f", defaults["Path"], "-f=path")
	f.Storage = flag.String("stype", defaults["Storage"], "-s=storage")
	f.DSN = flag.String("d", "", "-d=connection_string")
	f.VDB = flag.String("vdb", "", "-vdb=virtual_db_name")
	f.HTTPS = flag.Bool("s", false, "-s to enable a HTTPS connection")
	f.Cfg = flag.String("c", "", "-c=path/to/conf.json")
	f.Config = flag.String("config", "", "-config=path/to/conf.json")
	f.TrustedSubNetwork = flag.String("t", "", "-t=trusted_subnet")
}

// Config contains all the settings for configuring the application.
type Config struct {
	Host              string
	BaseURL           string
	TrustedSubNetwork *net.IPNet
	Key               []byte
	DBConfig          *repository.Config
	HTTPS             bool
}

// Modify modifies the config by the file provided.
func Modify(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("can't open %s: %v", filename, err)
	}
	defer file.Close()

	all, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("can't read %s: %v", filename, err)
	}

	var fCopy Flag
	err = json.Unmarshal(all, &fCopy)
	if err != nil {
		return fmt.Errorf("can't unmarshal %s: %v", filename, err)
	}

	// Reflection by the Flag structure and replacement of null attributes by config file provided.
	reflectionF := reflect.ValueOf(&f).Elem()
	reflectionFCopy := reflect.ValueOf(&fCopy).Elem()

	for i := 0; i < reflectionF.NumField(); i++ {
		field := reflectionF.Type().Field(i)
		fieldKind := reflectionF.Field(i)
		switch fieldKind.Kind() {
		case reflect.Ptr:
			elem := fieldKind.Elem()
			switch elem.Type().Kind() {
			case reflect.String:
				if val := defaults[field.Name]; (elem.String() == "" || val == elem.String()) &&
					reflectionFCopy.Field(i).Elem().IsValid() {

					elem.SetString(reflectionFCopy.Field(i).Elem().String())
				}
			case reflect.Bool:
				if !elem.Bool() {
					elem.SetBool(reflectionFCopy.Field(i).Elem().Bool())
				}
			}

		}
	}

	return nil
}

// New initializing the config for the application.
func New() *Config {
	flag.Parse()

	configFile := *f.Cfg
	if cfg, ok := os.LookupEnv("CONFIG"); configFile == "" && !ok {
		configFile = *f.Config
	} else if ok {
		configFile = cfg
	}

	if configFile != "" {
		err := Modify(configFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if addr, ok := os.LookupEnv("BASE_URL"); ok {
		f.BaseURL = &addr
	}

	if fsp, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		f.Path = &fsp
	}

	if sub, ok := os.LookupEnv("TRUSTED_SUBNET"); ok {
		f.TrustedSubNetwork = &sub
	}

	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		f.Host = &addr
	}

	if dsn, ok := os.LookupEnv("DATABASE_DSN"); ok {
		f.DSN = &dsn
	}

	if _, ok := os.LookupEnv("ENABLE_HTTPS"); ok {
		f.HTTPS = &ok
	}

	if *f.DSN == "" && storage.Type(*f.Storage) == defaultStorage && *f.VDB == "" {
		s := ""
		if *f.Path != defaultPath {
			s = string(filestorage.FileStorageType)
			f.Storage = &s
		} else {
			s = string(mapstorage.MapStorageType)
			f.Storage = &s
		}
	}

	var ddb *dockerdb.VDB
	var vdb = *f.VDB

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
				Name:  *f.Storage,
				Image: *f.Storage,
			},
		}

		var err error
		ddb, err = dockerdb.New(ctx, cfg)
		if err != nil {
			log.Fatal(err)
		}
		f.DSN = &ddb.ConnString
	}

	_, subnet, _ := net.ParseCIDR(*f.TrustedSubNetwork)

	var config = &Config{
		Host:    *f.Host,
		BaseURL: *f.BaseURL,
		Key:     []byte("CHANGE ME"),
		DBConfig: &repository.Config{
			DriverName:     storage.Type(*f.Storage),
			DataSourcePath: *f.Path,
			DataSourceCred: *f.DSN,
			VDB:            ddb,
			Name:           vdb,
		},
		HTTPS:             *f.HTTPS,
		TrustedSubNetwork: subnet,
	}

	return config
}
