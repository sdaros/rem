package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", handlers.AddReminder())
	err := http.ListenAndServe(":42888", mux)
	log.Fatal(err)

}
