package server

import (
	"github.com/shilin-anton/urlreducer/internal/app/config"
	"github.com/shilin-anton/urlreducer/internal/app/handlers"
	"log"
	"net/http"
)

type Storage interface {
	Add(short string, url string)
	Get(short string) (string, bool)
	FindByValue(url string) (string, bool)
	Size() int
}

type server struct {
	handler http.Handler
	storage Storage
}

func New(storage Storage) *server {
	handler := handlers.New(storage)
	S := &server{
		handler: handler,
		storage: storage,
	}
	return S
}

func (s server) Start() {
	log.Fatal(http.ListenAndServe(config.RunAddr, s.handler))
}
