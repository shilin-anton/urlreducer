package server

import (
	"github.com/shilin-anton/urlreducer/internal/app/handlers"
	"log"
	"net/http"
)

type Storage interface {
	Add(short string, url string)
	Get(short string) (string, bool)
}

type server struct {
	host    string
	handler http.Handler
	storage Storage
}

func New(storage Storage) *server {
	handler := handlers.New()
	S := &server{
		host:    ":8080",
		handler: handler,
		storage: storage,
	}
	return S
}

func (s server) Start() {
	log.Fatal(http.ListenAndServe(":8080", s.handler))
}
