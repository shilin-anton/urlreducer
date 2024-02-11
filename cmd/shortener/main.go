package main

import (
	"github.com/shilin-anton/urlreducer/internal/app/config"
	"github.com/shilin-anton/urlreducer/internal/app/server"
	"github.com/shilin-anton/urlreducer/internal/app/storage"
)

func main() {
	config.ParseFlags()
	
	myStorage := storage.New()
	myServer := server.New(myStorage)
	myServer.Start()
}
