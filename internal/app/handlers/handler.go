package handlers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/shilin-anton/urlreducer/internal/app/config"
	"github.com/shilin-anton/urlreducer/internal/app/gzip"
	"github.com/shilin-anton/urlreducer/internal/app/storage"
	"github.com/shilin-anton/urlreducer/internal/logger"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Storage interface {
	Add(short string, url string)
	Get(short string) (string, bool)
}

type Server struct {
	data    Storage
	handler http.Handler
}

// types for logger
type responseData struct {
	status int
	size   int
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Result string `json:"result"`
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.responseData.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(data)
	lrw.responseData.size += size
	return size, err
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func New() *Server {
	r := chi.NewRouter()

	r.Use(requestLoggerMiddleware)
	r.Use(responseLoggerMiddleware)

	s := &Server{
		data:    make(storage.Storage),
		handler: r,
	}
	r.Get("/{short}", gzipMiddleware(s.GetHandler))
	r.Post("/", gzipMiddleware(s.PostHandler))
	r.Post("/api/shorten", gzipMiddleware(s.PostShortenHandler))

	return s
}

func shortenURL(url string) string {
	// Решил использовать хэширование и первые символы результата, как короткую форму URL
	hash := md5.Sum([]byte(url))
	hashString := hex.EncodeToString(hash[:])
	shortURL := hashString[:8]
	return shortURL
}

func (s Server) PostHandler(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	url := string(body)
	short := shortenURL(url)

	s.data.Add(short, url)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.BaseAddr + "/" + short))
}

func (s Server) GetHandler(res http.ResponseWriter, req *http.Request) {
	short := chi.URLParam(req, "short")

	url, ok := s.data.Get(short)
	if !ok {
		http.NotFound(res, req)
		return
	}
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (s Server) PostShortenHandler(res http.ResponseWriter, req *http.Request) {
	var request shortenRequest
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if request.URL == "" {
		http.Error(res, "url must be passed", http.StatusUnprocessableEntity)
		return
	}

	short := shortenURL(request.URL)
	s.data.Add(short, request.URL)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	response := shortenResponse{
		Result: config.BaseAddr + "/" + short,
	}

	enc := json.NewEncoder(res)
	if err = enc.Encode(response); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.RequestLogger(r.RequestURI, r.Method, time.Since(start).String())
		next.ServeHTTP(w, r)
	})
}

func responseLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w, responseData: &responseData{}}
		next.ServeHTTP(lrw, r)
		logger.ResponseLogger(strconv.Itoa(lrw.responseData.status), strconv.Itoa(lrw.responseData.size))
	})
}

func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := gzip.NewCompressWriter(w)
			ow = cw
			defer cw.Close()
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := gzip.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}
		h.ServeHTTP(ow, r)
	}
}
