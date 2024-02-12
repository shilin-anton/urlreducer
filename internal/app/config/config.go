package config

import (
	"flag"
	"fmt"
	"os"
)

var RunAddr string
var BaseAddr string

const defaultRunURL = "localhost:8080"
const defaultBaseURL = "http://localhost:8080"

func ParseConfig() {
	flag.StringVar(&RunAddr, "a", defaultRunURL, "address and port to run server")
	flag.StringVar(&BaseAddr, "b", defaultBaseURL, "base URL before short link")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		RunAddr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		BaseAddr = envBaseAddr
	}
	fmt.Printf("Server is running on %s\nBase URL is %s\n", RunAddr, BaseAddr)
}
