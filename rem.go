package main

import (
	"log"
	"net/http"
	"os"
	"sync"
)

type App struct {
	registry map[string]interface{}
	sync.Mutex
}

func main() {
	app := new(App)
	app.registry = make(map[string]interface{})
	getFromEnvOrSetDefaults(app)
	mux := http.NewServeMux()
	mux.Handle("/", AddReminder(app))
	mux.Handle("/create", CreateForm(app))
	log.Printf("Serving %v (version: %v) on %v%v", app.Lookup("name"),
		app.Lookup("version"), app.Lookup("port"), app.Lookup("path"))
	err := http.ListenAndServe(app.Lookup("port").(string), mux)
	log.Fatal(err)

}

func (a *App) Register(k string, v interface{}) {
	if a == nil {
		return
	}

	a.Lock()
	defer a.Unlock()
	a.registry[k] = v
}

func (a *App) Lookup(k string) interface{} {
	if a == nil {
		return nil
	}

	a.Lock()
	defer a.Unlock()
	return a.registry[k]
}

func getFromEnvOrSetDefaults(app *App) {
	app.Register("version", "0.1.0")
	app.Register("name", "rem")
	home := os.Getenv("HOME")
	if home == "" {
		log.Fatalln("HOME Environment variable not set.")
	}
	app.Register("home", home)
	port := os.Getenv("REM_PORT")
	if port == "" {
		port = "42888"
	}
	app.Register("port", ":"+port)
	path := os.Getenv("REM_PATH")
	if path == "" {
		path = "/rem/"
	}
	app.Register("path", path)
}
