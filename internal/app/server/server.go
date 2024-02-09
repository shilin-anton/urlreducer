package server

import (
	"net/http"

	"github.com/shilin-anton/urlreducer/internal/app/handlers"
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
		host:    "localhost:8080",
		handler: handler,
		storage: storage,
	}
	return S
}

func (s server) Start() error {
	err := http.ListenAndServe(s.host, s.handler)
	if err != nil {
		return err
	}
	return nil
}
