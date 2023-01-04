package main

import (
	"compress/gzip"
	"flag"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"url-shortener/config"
	handlers "url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/routes"
)

var (
	host    *string
	baseURL *string
	path    *string
)

func init() {
	host = flag.String("a", "localhost:8080", "-a=host")
	baseURL = flag.String("b", "http://localhost:8080/", "-b=URL")
	path = flag.String("f", "urlshortener.txt", "-f=path")
}

func main() {
	flag.Parse()

	cfg := config.New(*baseURL, *path)

	storage, err := repository.New(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}

	router := gin.Default()
	handler := handlers.NewHandler(storage)
	public := router.Group("/")
	routes.PublicRoutes(public, *handler)

	//router.Use(func(c *gin.Context) {
	//	if str := c.Request.Header.Get("Accept-Encoding"); str != "" && strings.Contains(str, "gzip") {
	//		c.Writer.Header().Set("Content-Encoding", "gzip")
	//		gzip.Gzip(gzip.BestSpeed)
	//	}
	//})

	//router.Use(gzip.Gzip(gzip.BestSpeed))

	serverAddress := *host
	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		serverAddress = addr
	}

	http.ListenAndServe(serverAddress, gzipHandle(router))
	//router.Run(serverAddress)
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
