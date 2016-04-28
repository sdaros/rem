package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	version = "0.3.0"
)

type (
	App struct {
		registry map[string]interface{}
		sync.Mutex
	}
	Config struct {
		DocumentRoot     string
		Path             string
		Port             string
		ReminderTemplate string
		RemScript        string
	}
)

func main() {
	app := new(App)
	app.registry = make(map[string]interface{})
	loadConfigurationFile(app)
	config := app.Lookup("config").(*Config)
	mux := http.NewServeMux()
	mux.Handle("/", CreateReminder(app))
	log.Printf("Serving rem (version: %v) on %v%v",
		version, config.Port, config.Path)
	err := http.ListenAndServe(config.Port, mux)
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

func loadConfigurationFile(app *App) {
	file, err := os.Open("./rem.conf")
	if err != nil {
		log.Fatal("error: Configuration file not found!")
	}
	decoder := json.NewDecoder(file)
	config := new(Config)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("error parsing configuration file: %v", err)
	}
	app.Register("config", config)
}
