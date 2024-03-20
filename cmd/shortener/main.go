package main

import (
	"github.com/shilin-anton/urlreducer/internal/app/config"
	"github.com/shilin-anton/urlreducer/internal/app/server"
	"github.com/shilin-anton/urlreducer/internal/app/storage"
	"github.com/shilin-anton/urlreducer/internal/logger"
	"log"
)

func main() {
	config.ParseConfig()
	err := logger.Initialize(config.LogLevel)
	if err != nil {
		log.Fatal("Error initializing logger: %v\n", err)
	}

	myStorage := storage.New()
	myServer := server.New(myStorage)
	myServer.Start()
}
