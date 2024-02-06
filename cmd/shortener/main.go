package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Storage struct {
	storage map[string]string
}

func (s *Storage) Add(short string, url string) {
	s.storage[short] = url
}

func (s *Storage) Get(short string) (string, bool) {
	url, ok := s.storage[short]
	return url, ok
}

func CreateStorage() *Storage {
	return &Storage{storage: make(map[string]string)}
}

var urlStorage = CreateStorage()

func shortenURL(url string) string {
	// Решил использовать хэширование и первые символы результата, как короткую форму URL
	hash := md5.Sum([]byte(url))
	hashString := hex.EncodeToString(hash[:])
	shortURL := hashString[:8]
	return shortURL
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, defineHandler)

	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func defineHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && req.URL.Path == "/" {
		postHandler(res, req)
	} else if req.Method == http.MethodGet && req.URL.Path != "/" {
		getHandler(res, req)
	} else {
		invalidRequestHandler(res, req)
	}
}

func postHandler(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	url := string(body)
	short := shortenURL(url)

	urlStorage.Add(short, url)

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("Content-Type", "text/plain")
	res.Write([]byte(short))
}

func getHandler(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/"):]
	if id == "" || strings.Contains(id, "/") {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	url, ok := urlStorage.Get(id)
	if !ok {
		http.NotFound(res, req)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Write(nil)
}

func invalidRequestHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Bad Request", http.StatusBadRequest)
}
