package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/go-chi/chi/v5"
	"github.com/shilin-anton/urlreducer/internal/app/config"
	"github.com/shilin-anton/urlreducer/internal/app/storage"
	"io"
	"net/http"
)

type Storage interface {
	Add(short string, url string)
	Get(short string) (string, bool)
}

type Server struct {
	data    Storage
	handler http.Handler
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func New() *Server {
	r := chi.NewRouter()

	S := &Server{
		data:    make(storage.Storage),
		handler: r,
	}
	r.Get("/{short}", S.GetHandler)
	r.Post("/", S.PostHandler)

	return S
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
	res.Write([]byte(config.BaseAddr + short))
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
