package main

import (
	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/fileutil"
	"github.com/KillReall666/yaproject/internal/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KillReall666/yaproject/internal/handlers/get"
	update2 "github.com/KillReall666/yaproject/internal/handlers/update"
	"github.com/KillReall666/yaproject/internal/service"
	"github.com/KillReall666/yaproject/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetHandler(t *testing.T) {
	type response struct {
		requestCode int
		contentType string
	}
	tests := []struct {
		name     string
		method   string
		response response
		url      string
	}{
		{
			name:   "too short url 1",
			method: http.MethodGet,
			url:    "/value/counter",
			response: response{
				requestCode: http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:   "too short url 2",
			method: http.MethodGet,
			url:    "/value",
			response: response{
				requestCode: http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:   "unknown type of metric",
			method: http.MethodGet,
			url:    "/value/unknown_type/Alloc/",
			response: response{
				requestCode: http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			cfg := config.RunFileIo{}
			log, _ := logger.InitLogger()
			store := storage.NewMemStorage()
			fileio := fileutil.NewFileIo(cfg, store, log)
			serv := service.NewService(store, log, fileio)
			get := get.NewGetHandler(serv)

			get.GetMetrics(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.response.requestCode, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.response.contentType)
		})
	}
}

func TestPostHandler(t *testing.T) {
	type response struct {
		requestCode int
		contentType string
	}
	tests := []struct {
		name     string
		method   string
		response response
		url      string
	}{
		{
			name:   "too short url 1",
			method: http.MethodPost,
			url:    "/update/gauge/Alloc",
			response: response{
				requestCode: http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:   "too short url 2",
			method: http.MethodPost,
			url:    "/update/counter",
			response: response{
				requestCode: http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:   "too short url 3",
			method: http.MethodPost,
			url:    "/update",
			response: response{
				requestCode: http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:   "unknown type of metric", // can crash if in PostHandler called h.service.MetricsPrint()
			method: http.MethodPost,
			url:    "/update/unknown_type/Alloc/",
			response: response{
				requestCode: http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			cfg := config.RunFileIo{}
			log, _ := logger.InitLogger()
			store := storage.NewMemStorage()
			fileio := fileutil.NewFileIo(cfg, store, log)
			serv := service.NewService(store, log, fileio)
			update := update2.NewUpdateHandler(serv)

			update.UpdateMetrics(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.response.requestCode, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.response.contentType)
		})
	}
}
