package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
)

const (
	version = "0.6.0"
)

type App struct {
	Notification
	Domain string
	Path   string
	Port   string
}

func main() {
	app := new(App)
	app.Notification = *&Notification{}
	app.loadConfigurationFile()
	mux := http.NewServeMux()
	mux.Handle("/", CreateReminder(app))
	log.Printf("Serving rem (version: %v) on %v/%v",
		version, ":"+app.Port, app.Path)
	err := http.ListenAndServe(":"+app.Port, mux)
	log.Fatal(err)
}

func (self *App) loadConfigurationFile() {
	var config = flag.String("config", "rem.conf",
		"configuration file to use")
	flag.Parse()
	configFile, err := os.Open(*config)
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
