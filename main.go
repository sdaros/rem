package main

import (
	"log"
	"net/http"

	"github.com/sdaros/rem/handlers"
)

const (
	// buildMetadata is replaced when package is built using -ldflags -X
	// ex: go build -ldflags "-X main.buildMetadata=+`git rev-parse --short HEAD`"
	buildMetadata = ""
	version       = "0.1.0"
	name          = "rem"
)

func main() {
	port := ":42888"
	mux := http.NewServeMux()
	mux.Handle("/", AddReminder())
	log.Printf("Serving %v (version: %v) on %v", name, Version(), port)
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)

}

func Version() string { return version + buildMetadata }
