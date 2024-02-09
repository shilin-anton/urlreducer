package main

import (
	"fmt"
	"github.com/shilin-anton/urlreducer/internal/app/server"
	"github.com/shilin-anton/urlreducer/internal/app/storage"
)

func main() {
	myStorage := storage.New()
	myServer := server.New(myStorage)
	err := myServer.Start()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
