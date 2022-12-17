package handler

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
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
				r.EXPECT().GetLink("/IVI").Return(
					"http://zrnzruvv7qfdy.ru/hlc65i", nil)
			},
			expectedStatusCode:   307,
			expectedResponseHead: `http://zrnzruvv7qfdy.ru/hlc65i`,
		},
		{
			name:   "Err",
			target: "/IVI1",
			mockBehavior: func(r *service_mocks.MockGetLink) {
				r.EXPECT().GetLink("/IVI1").Return(
					"", errors.New("Bad url"))
			},
			expectedStatusCode:   400,
			expectedResponseHead: ``,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockGetLink(c)
			test.mockBehavior(repo)

			services := &service.Service{GetLink: repo}
			handler := Handler{services}

			req := httptest.NewRequest("GET", test.target,
				bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(handler.GetLinkHandler)
			// запускаем сервер
			h.ServeHTTP(w, req)
			//res := w.Result()

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
					"", errors.New("UNIQUE constraint failed: urls.long"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: "UNIQUE constraint failed: urls.long\n",
		},
		{
			name:      "server error",
			inputBody: "q",
			mockBehavior: func(r *service_mocks.MockCreateLink) {
				r.EXPECT().CreateLink("q").Return(
					"", errors.New(""))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "недопустимый URL\n\n",
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
			h := http.HandlerFunc(handler.CreateLinkHandler)
			// запускаем сервер
			h.ServeHTTP(w, req)
			//res := w.Result()

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

//type dbMock struct {
//}
//
//func (dm dbMock) Close() error {
//	fmt.Println("Closed successfully")
//	return nil
//}
//
//func (dm dbMock) Exec(query string, args ...any) (sql.Result, error) {
//	return nil, nil
//}
//
//func (dm dbMock) QueryRow(query string, args ...any) *sql.Row {
//	return &sql.Row{}
//}

//
//func TestGetHandler(t *testing.T) {
//	type want struct {
//		Status      int
//		ContentType string
//		Response    string
//	}
//	tests := []struct {
//		name    string
//		storage storage.Repositories
//		want    want
//	}{
//		{
//			"404 test #1",
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			id := getTestID()
//			request := httptest.NewRequest(http.MethodGet, "/"+id, nil)
//		})
//	}
//}
