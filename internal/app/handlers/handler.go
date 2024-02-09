package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/shilin-anton/urlreducer/internal/app/storage"
	"io"
	"net/http"
	"strings"
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
	mux := http.NewServeMux()
	S := &Server{
		data:    make(storage.Storage),
		handler: mux,
	}
	mux.HandleFunc("/", S.DefineHandler)
	return S
}

func shortenURL(url string) string {
	// Решил использовать хэширование и первые символы результата, как короткую форму URL
	hash := md5.Sum([]byte(url))
	hashString := hex.EncodeToString(hash[:])
	shortURL := hashString[:8]
	return shortURL
}

func (s Server) DefineHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && req.URL.Path == "/" {
		s.PostHandler(res, req)
	} else if req.Method == http.MethodGet && req.URL.Path != "/" {
		s.GetHandler(res, req)
	} else {
		s.InvalidRequestHandler(res, req)
	}
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
	res.Write([]byte("http://localhost:8080/" + short))
}

func (s Server) GetHandler(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/"):]
	if id == "" || strings.Contains(id, "/") {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	url, ok := s.data.Get(id)
	if !ok {
		http.NotFound(res, req)
		return
	}
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (s Server) InvalidRequestHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Bad Request", http.StatusBadRequest)
}
