package config

import (
	"flag"
	"os"
)

var RunAddr string
var BaseAddr string
var LogLevel string

const defaultRunURL = "localhost:8080"
const defaultBaseURL = "http://localhost:8080"
const defaultLogLevel = "info"

func ParseConfig() {
	flag.StringVar(&RunAddr, "a", defaultRunURL, "address and port to run server")
	flag.StringVar(&BaseAddr, "b", defaultBaseURL, "base URL before short link")
	flag.StringVar(&LogLevel, "l", defaultLogLevel, "log level")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		RunAddr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		BaseAddr = envBaseAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		LogLevel = envLogLevel
	}
}
