package main

import (
	"github.com/shilin-anton/urlreducer/internal/app/config"
	file_manager "github.com/shilin-anton/urlreducer/internal/app/file-manager"
	"github.com/shilin-anton/urlreducer/internal/app/server"
	"github.com/shilin-anton/urlreducer/internal/app/storage"
	"github.com/shilin-anton/urlreducer/internal/logger"
	"log"
)

func main() {
	config.ParseConfig()
	err := logger.Initialize(config.LogLevel)
	if err != nil {
		log.Fatal("Error initializing logger")
	}

	myStorage := storage.New()
	file_manager.ReadFromFile(&myStorage)
	myServer := server.New(myStorage)
	myServer.Start()
}
