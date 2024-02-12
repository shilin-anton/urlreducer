package main

import (
	"github.com/shilin-anton/urlreducer/internal/app/config"
	"github.com/shilin-anton/urlreducer/internal/app/server"
	"github.com/shilin-anton/urlreducer/internal/app/storage"
)

func main() {
	config.ParseConfig()

	myStorage := storage.New()
	myServer := server.New(myStorage)
	myServer.Start()
}
