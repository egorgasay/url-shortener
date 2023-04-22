package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"url-shortener/config"
	"url-shortener/internal/repository"
	"url-shortener/internal/usecase"
)

func TestHandler_GetLinkHandler(t *testing.T) {
	tests := []struct {
		name                 string
		target               string
		expectedStatusCode   int
		expectedResponseHead string
	}{
		{
			name:                 "Ok",
			target:               "/zE",
			expectedStatusCode:   307,
			expectedResponseHead: "http://zrnzruvv7qfdy.ru/hlc65i",
		},
		{
			name:                 "Err",
			target:               "/IVI1",
			expectedStatusCode:   400,
			expectedResponseHead: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &repository.Config{
				DriverName:     "map",
				DataSourcePath: "test",
			}
			repo, err := repository.New(cfg)
			if err != nil {
				t.Fatal(err)
			}

			logic := usecase.New(repo)

			repo.AddLink("http://zrnzruvv7qfdy.ru/hlc65i", "zE", "df")

			conf := &config.Config{Host: "127.0.0.1", DBConfig: cfg}
			handler := Handler{conf: conf, logic: logic}

			req := httptest.NewRequest("GET", test.target,
				nil)
			w := httptest.NewRecorder()

			router := gin.Default()
			router.GET("/:id", handler.GetLinkHandler)

			router.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Header().Get("Location"), test.expectedResponseHead)
		})
	}
}

func TestHandler_CreateLinkHandler(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Ok",
			inputBody:            "vk.com/gasayminajj",
			expectedStatusCode:   201,
			expectedResponseBody: `zE`,
		},
		{
			name:                 "server error",
			inputBody:            "q",
			expectedStatusCode:   500,
			expectedResponseBody: "",
		},
		{
			name:                 "Exists",
			inputBody:            "vk.com/gasayminajj",
			expectedStatusCode:   409,
			expectedResponseBody: `rx`,
		},
	}

	cfg := &repository.Config{
		DriverName: "map",
	}

	repo, err := repository.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	logic := usecase.New(repo)
	conf := &config.Config{Host: "127.0.0.1", DBConfig: cfg}
	handler := Handler{conf: conf, logic: logic}

	// определяем хендлер
	router := gin.Default()
	router.POST("/", handler.CreateLinkHandler)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/",
				bytes.NewBufferString(test.inputBody))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
			repo.AddLink("vk.com/gasayminajj", "rx", "80f53850d88d388b2a5fb1a057a1867ee70d37b1c2439ede79c43ef3c802e4b8-31363832313934313833373336353432343636")
		})
	}
}

func TestHandler_APICreateLinkHandler(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Ok",
			inputBody:            `{"url":"vk.com/gasayminajj"}`,
			expectedStatusCode:   201,
			expectedResponseBody: `{"result":"zE"}`,
		},
		{
			name:                 "server error",
			inputBody:            "q",
			expectedStatusCode:   500,
			expectedResponseBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &repository.Config{
				DriverName:     "map",
				DataSourcePath: "test",
			}
			repo, err := repository.New(cfg)
			if err != nil {
				t.Fatal(err)
			}

			logic := usecase.New(repo)

			conf := &config.Config{Host: "127.0.0.1", DBConfig: cfg}
			handler := Handler{conf: conf, logic: logic}

			req := httptest.NewRequest("POST", "/api/shorten",
				bytes.NewBufferString(test.inputBody))
			w := httptest.NewRecorder()

			router := gin.Default()
			router.Use(handler.APICreateLinkHandler)

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_Ping(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
	}{
		{
			name:               "Ok",
			expectedStatusCode: 200,
		},
	}

	cfg := &repository.Config{
		DriverName:     "map",
		DataSourcePath: "test",
	}
	repo, err := repository.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	logic := usecase.New(repo)

	conf := &config.Config{Host: "127.0.0.1", DBConfig: cfg}
	handler := Handler{conf: conf, logic: logic}

	router := gin.Default()
	router.Use(handler.Ping)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/",
				bytes.NewBufferString(""))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_BatchHandler(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Ok",
			inputBody:            `[{"correlation_id": "1", "original_url": "vk.com/gasayminajj"}]`,
			expectedStatusCode:   201,
			expectedResponseBody: "[\n    {\n        \"correlation_id\": \"1\",\n        \"short_url\": \"1\"\n    }\n]",
		},
		{
			name:                 "Bad JSON",
			inputBody:            `[{"correlati", "original_url": "vk.com/gasayminajj"}]`,
			expectedStatusCode:   400,
			expectedResponseBody: "",
		},
	}

	cfg := &repository.Config{
		DriverName:     "map",
		DataSourcePath: "test",
	}
	repo, err := repository.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	logic := usecase.New(repo)

	conf := &config.Config{Host: "127.0.0.1", DBConfig: cfg}
	handler := Handler{conf: conf, logic: logic}

	router := gin.Default()
	router.Use(handler.BatchHandler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/shorten/batch",
				bytes.NewBufferString(tt.inputBody))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_APIDeleteLinksHandler(t *testing.T) {
	cfg := &repository.Config{
		DriverName:     "map",
		DataSourcePath: "test",
	}
	repo, err := repository.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	logic := usecase.New(repo)

	conf := &config.Config{Host: "127.0.0.1", DBConfig: cfg}
	handler := Handler{conf: conf, logic: logic}

	inputBody := `[ "zE" ]`

	_, err = repo.AddLink("http://zrnzruvv7qfdy.ru/hlc65i", "zE", "80f53850d88d388b2a5fb1a057a1867ee70d37b1c2439ede79c43ef3c802e4b8-31363832313934313833373336353432343636")
	if err != nil {
		t.Error(err)
	}

	router := gin.Default()
	router.GET("/:id", handler.GetLinkHandler)
	router.DELETE("/api/user/urls", handler.APIDeleteLinksHandler)

	req := httptest.NewRequest("DELETE", "/api/user/urls",
		bytes.NewBufferString(inputBody))
	req.Header.Set("Authorization", "80f53850d88d388b2a5fb1a057a1867ee70d37b1c2439ede79c43ef3c802e4b8-31363832313934313833373336353432343636")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	time.Sleep(1 * time.Second)

	req = httptest.NewRequest("GET", "/zE", nil)

	w = httptest.NewRecorder()

	t.Log(w.Body.String())
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusGone, w.Code)
}
