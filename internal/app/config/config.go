package config

import (
	"flag"
)

var RunAddr string
var BaseAddr string

func ParseFlags() {
	flag.StringVar(&RunAddr, "a", "localhost:8080/", "address and port to run server")
	flag.StringVar(&BaseAddr, "b", "http://localhost:8080/", "base URL before short link")
	flag.Parse()
}
