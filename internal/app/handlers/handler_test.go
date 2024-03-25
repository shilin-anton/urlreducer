package handlers

import (
	"github.com/shilin-anton/urlreducer/internal/app/config"
	"github.com/stretchr/testify/assert"
	"io"
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

func TestServer_PostShortenHandler(t *testing.T) {
	srv := New()
	config.BaseAddr = "http://localhost:8080"

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "method_post_without_body",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "method_post_unsupported_type",
			method:       http.MethodPost,
			body:         `{"url": ""}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: "",
		},
		{
			name:         "method_post_success",
			method:       http.MethodPost,
			body:         `{"url": "https://yandex.ru"}`,
			expectedCode: http.StatusCreated,
			expectedBody: `{"result": "http://localhost:8080/e9db20b2"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/shorten", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			if len(tc.body) > 0 {
				req.Header.Set("Content-Type", "application/json")
			}

			srv.handler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, tc.expectedCode)
			}
			if tc.expectedBody != "" {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				bodyString := string(body)
				assert.JSONEq(t, tc.expectedBody, bodyString)
			}
		})
	}
}
