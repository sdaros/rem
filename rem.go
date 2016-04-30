package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
)

const (
	version = "0.3.0"
)

type App struct {
	DocumentRoot string
	Path         string
	Port         string
	RemScript    string
	ApiToken     string
	ApiUser      string
}

func main() {
	app := new(App)
	loadConfigurationFile(app)
	wantsInit := flag.Bool("init", false, "initialise REM on Uberspace")
	flag.Parse()
	if *wantsInit {
		initScriptToStdout(app)
		return
	}
	mux := http.NewServeMux()
	mux.Handle("/", CreateReminder(app))
	log.Printf("Serving rem (version: %v) on %v/%v",
		version, app.Port, app.Path)
	err := http.ListenAndServe(app.Port, mux)
	log.Fatal(err)

}

func loadConfigurationFile(app *App) {
	home := os.Getenv("HOME")
	configFile, err := os.Open(home + "/.config/rem/rem.conf")
	if err != nil {
		log.Fatal("error: Configuration file not found!")
	}
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&app)
	if err != nil {
		log.Fatalf("error parsing configuration file: %v", err)
	}
}

func initScriptToStdout(app *App) {
	initScript := app.DocumentRoot + "/" + app.Path + "/init_script.template"
	t, err := template.ParseFiles(initScript)
	if err != nil {
		log.Fatalf("error parsing install.template file: %v", err)
	}
	t.Execute(os.Stdout, app)
}
