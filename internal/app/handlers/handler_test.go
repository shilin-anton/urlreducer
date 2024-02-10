package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostHandler(t *testing.T) {
	srv := New()

	tests := []struct {
		name           string
		method         string
		url            string
		requestBody    string
		wantStatusCode int
	}{
		{
			name:           "Valid POST request",
			method:         http.MethodPost,
			url:            "/",
			requestBody:    "http://example.com",
			wantStatusCode: http.StatusCreated,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.url, strings.NewReader(test.requestBody))
			w := httptest.NewRecorder()

			srv.handler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != test.wantStatusCode {
				t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, test.wantStatusCode)
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	srv := New()
	srv.data.Add("test_short", "https://smth.ru")

	tests := []struct {
		name               string
		method             string
		url                string
		wantStatusCode     int
		wantLocationHeader string
	}{
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()

			srv.handler.ServeHTTP(w, req)

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
