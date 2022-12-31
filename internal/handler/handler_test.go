package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/service"
	service_mocks "url-shortener/internal/service/mocks"
)

func TestHandler_GetLinkHandler(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockGetLink)

	tests := []struct {
		name                 string
		target               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseHead string
	}{
		{
			name:   "Ok",
			target: "/IVI",
			mockBehavior: func(r *service_mocks.MockGetLink) {
				r.EXPECT().GetLink("IVI").Return(
					"http://zrnzruvv7qfdy.ru/hlc65i", nil).AnyTimes()
			},
			expectedStatusCode:   307,
			expectedResponseHead: "http://zrnzruvv7qfdy.ru/hlc65i",
		},
		{
			name:   "Err",
			target: "/IVI1",
			mockBehavior: func(r *service_mocks.MockGetLink) {
				r.EXPECT().GetLink("IVI1").Return(
					"", errors.New("bad url")).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseHead: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockGetLink(c)
			test.mockBehavior(repo)

			services := &service.Service{GetLink: repo}
			handler := Handler{services}

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
	type mockBehavior func(r *service_mocks.MockCreateLink)

	tests := []struct {
		name                 string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: "vk.com/gasayminajj",
			mockBehavior: func(r *service_mocks.MockCreateLink) {
				r.EXPECT().CreateLink("vk.com/gasayminajj").Return(
					"BEh6", nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `http://127.0.0.1:8080/BEh6`,
		},
		{
			name:      "already exists",
			inputBody: "http://zrnzqddy.ru/hlc65i",
			mockBehavior: func(r *service_mocks.MockCreateLink) {
				r.EXPECT().CreateLink("http://zrnzqddy.ru/hlc65i").Return(
					"", gin.Error{Err: errors.New("URL уже существует")})
			},
			expectedStatusCode:   500,
			expectedResponseBody: "",
		},
		{
			name:      "server error",
			inputBody: "q",
			mockBehavior: func(r *service_mocks.MockCreateLink) {
				r.EXPECT().CreateLink("q").Return(
					"", gin.Error{Err: errors.New("недопустимый URL")}).
					AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockCreateLink(c)
			test.mockBehavior(repo)

			services := &service.Service{CreateLink: repo}
			handler := Handler{services}

			req := httptest.NewRequest("POST", "/",
				bytes.NewBufferString(test.inputBody))
			w := httptest.NewRecorder()
			// определяем хендлер
			router := gin.Default()
			router.Use(handler.CreateLinkHandler)

			router.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
