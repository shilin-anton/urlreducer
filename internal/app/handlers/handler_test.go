package handlers

import (
	"github.com/shilin-anton/urlreducer/internal/app/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefineHandler(t *testing.T) {
	store := storage.New()
	store.Add("test_short", "https://smth.ru")

	srv := &Server{
		data:    store,
		handler: nil,
	}

	tests := []struct {
		name               string
		method             string
		url                string
		requestBody        string
		wantStatusCode     int
		wantLocationHeader string
	}{
		{
			name:           "Valid POST request",
			method:         http.MethodPost,
			url:            "/",
			requestBody:    "http://example.com",
			wantStatusCode: http.StatusCreated,
		},
		{
			name:               "Valid GET request with existing short link",
			method:             http.MethodGet,
			url:                "/test_short",
			wantStatusCode:     http.StatusTemporaryRedirect,
			wantLocationHeader: "https://smth.ru",
		},
		{
			name:           "Invalid GET request with non-existing short link",
			method:         http.MethodGet,
			url:            "/non_existing_short_link",
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "Invalid request",
			method:         http.MethodDelete,
			url:            "/",
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.url, strings.NewReader(test.requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(srv.DefineHandler)
			h(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != test.wantStatusCode {
				t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, test.wantStatusCode)
			}

			if test.wantLocationHeader != "" {
				location := resp.Header.Get("Location")
				if location != test.wantLocationHeader {
					t.Errorf("unexpected Location header: got %s, want %s", location, test.wantLocationHeader)
				}
			}
		})
	}
}
