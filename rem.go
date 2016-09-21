package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const (
	version = "0.6.0"
)

type App struct {
	Notification
	Port string
	Url  string
}

func main() {
	app := new(App)
	app.Notification = *&Notification{}
	app.loadConfigurationFile()
	mux := http.NewServeMux()
	mux.Handle("/", CreateReminder(app))
	log.Printf("Serving rem (version: %v) on %v",
		version, app.Port)
	err := http.ListenAndServe(app.Port, mux)
	log.Fatal(err)
}

func (self *App) loadConfigurationFile() {
	homeDir := os.Getenv("HOME")
	configFile, err := os.Open(homeDir + "/.config/rem/rem.conf")
	die("error: unable to find configuration file! %v", err)
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&self)
	die("error: unable to parse configuration file: %v", err)
}

func die(format string, err error) {
	if err != nil {
		log.Fatalf(format, err)
	}
}
