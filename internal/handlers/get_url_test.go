package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUrl(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string
		want   []string
	}{
		{
			name:   "first post-request test with value",
			method: http.MethodPost,
			url:    "/html/gauge/alloc/666",
			want:   []string{"html", "gauge", "alloc", "666"},
		},
		{
			name:   "second post-request test without value",
			method: http.MethodPost,
			url:    "/html/gauge/alloc",
			want:   []string{"html", "gauge", "alloc"},
		},
		{
			name:   "first get-request test with value",
			method: http.MethodGet,
			url:    "/value/gauge/alloc/666",
			want:   []string{"value", "gauge", "alloc", "666"},
		},
		{
			name:   "second get-request test without value",
			method: http.MethodGet,
			url:    "/value/gauge/alloc",
			want:   []string{"value", "gauge", "alloc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.url, nil)
			assert.Equal(t, GetURL(r), tt.want)
		})
	}
}
