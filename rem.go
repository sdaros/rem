package main

import (
	"log"
	"net/http"

	"github.com/sdaros/rem/handlers"
)

var (
	version       = "0.1.0"
	buildMetadata string
	name          = "rem"
)

func main() {
	port := ":42888"
	mux := http.NewServeMux()
	mux.Handle("/", handlers.AddReminder())
	log.Printf("Serving %v (version: %v) on %v", name, Version(), port)
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)

}

func Version() string { return version + buildMetadata }
