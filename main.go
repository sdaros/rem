package main

import (
	"log"
	"net/http"
)

func main() {
	port := ":42888"
	mux := http.NewServeMux()
	mux.Handle("/", AddReminder())
	log.Printf("Serving on %v", port)
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)

}
