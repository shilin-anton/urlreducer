package config

import (
	"flag"
	"os"
)

var RunAddr string
var BaseAddr string
var LogLevel string
var FilePath string

const defaultRunURL = "localhost:8080"
const defaultBaseURL = "http://localhost:8080"
const defaultLogLevel = "info"
const defaultFilePath = "/tmp/short-url-db.json"

func ParseConfig() {
	flag.StringVar(&RunAddr, "a", defaultRunURL, "address and port to run server")
	flag.StringVar(&BaseAddr, "b", defaultBaseURL, "base URL before short link")
	flag.StringVar(&LogLevel, "l", defaultLogLevel, "log level")
	flag.StringVar(&FilePath, "f", defaultFilePath, "file storage path")

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
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		FilePath = envFilePath
	}
}
