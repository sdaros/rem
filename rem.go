package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const (
	version = "0.5.0"
)

type App struct {
	*Notification
	DocumentRoot string
	Domain       string
	Path         string
	Port         string
}

func main() {
	app := new(App)
	app.Notification = new(Notification)
	loadConfigurationFile(app)
	mux := http.NewServeMux()
	mux.Handle("/", CreateReminder(app))
	log.Printf("Serving rem (version: %v) on %v/%v",
		version, app.Port, app.Path)
	err := http.ListenAndServe(app.Port, mux)
	log.Fatal(err)
}

func loadConfigurationFile(app *App) {
	homeDir := os.Getenv("HOME")
	configFile, err := os.Open(homeDir + "/.config/rem/rem.conf")
	die("error: unable to find configuration file! %v", err)
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&app)
	die("error: unable to parse configuration file: %v", err)
}

func die(format string, err error) {
	if err != nil {
		log.Fatalf(format, err)
	}
}
