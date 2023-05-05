package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/config"
	"url-shortener/internal/repository"
	"url-shortener/internal/usecase"
)

func Benchmark_CreateAndGetLink(b *testing.B) {
	s := 1000
	cfg := &repository.Config{
		DriverName:     "test",
		DataSourcePath: "testsqlite3",
	}
	repo, err := repository.New(cfg)
	if err != nil {
		b.Fatal(err)
	}

	logic := usecase.New(repo)

	conf := &config.Config{Host: "127.0.0.1:8080", DBConfig: cfg}
	handler := NewHandler(conf, logic)

	router := gin.Default()
	router.Use(handler.CreateLinkHandler)
	req := httptest.NewRequest("POST", "/", nil)

	b.ResetTimer()
	b.Run("Create", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf("%d", s)))
			b.StartTimer()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			b.StopTimer()
			log.Println()
			s++
		}
	})

	uindex := 0
	router = gin.Default()
	router.GET("/:id", handler.GetLinkHandler)
	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req = httptest.NewRequest("GET", "/Xz",
				nil)
			w := httptest.NewRecorder()

			b.StartTimer()
			router.ServeHTTP(w, req)
			b.StopTimer()
			uindex++
		}
	})
}
