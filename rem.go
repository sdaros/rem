package main

import (
	"log"
	"net/http"
	"sync"
)

var (
	version       = "0.1.0"
	buildMetadata string
	name          = "rem"
)

type Registrar struct {
	values map[string]interface{}
	sync.Mutex
}

func main() {
	port := ":42888"
	mux := http.NewServeMux()
	mux.Handle("/", AddReminder())
	log.Printf("Serving %v (version: %v) on %v", name, Version(), port)
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)

}

func (r *Registrar) Register(k string, v interface{}) {
	if r == nil {
		return
	}

	r.Lock()
	defer r.Unlock()
	r.values[k] = v
}

func (r *Registrar) Lookup(k string) interface{} {
	if r == nil {
		return nil
	}

	r.Lock()
	defer r.Unlock()
	return r.values[k]
}

func Version() string { return version + buildMetadata }
