package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type templateData struct {
	ReminderSuccess string
	Path            string
	DocumentRoot    string
	RemScript       string
}

func CreateReminder(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := initialiseTemplateData(w, app)
		switch r.Method {
		case "GET":
			renderTemplate(w, "create", data)
			return
		case "POST":
			submitReminder(w, r, data)
			return
		}
	})
}

func initialiseTemplateData(w http.ResponseWriter, app *App) *templateData {
	path := app.Lookup("path").(string)
	documentRoot := app.Lookup("documentRoot").(string)
	reminderSuccess := ""
	remScript := app.Lookup("home").(string) + "/bin/rem_script"
	return &templateData{reminderSuccess, path, documentRoot, remScript}
}

func submitReminder(w http.ResponseWriter, r *http.Request, data *templateData) {
	thenDay := r.FormValue("day")
	thenTime := r.FormValue("time")
	message := r.FormValue("message")
	delay, err := timeToSleepFor(thenDay, thenTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data.ReminderSuccess = "Your reminder has been added, thank you!"
	renderTemplate(w, "create", data)
	log.Printf("The reminder '%v' will be sent to you at %v", message, time.Now().Add(delay))
	go func(time.Duration, *templateData, string) {
		select {
		case <-time.After(delay):
			execute(data, message)
		}
	}(delay, data, message)
}

func execute(data *templateData, msg string) {
	cmd := exec.Command(data.RemScript, msg)
	var cmdResult bytes.Buffer
	cmd.Stdout = &cmdResult
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func timeToSleepFor(thenDay, thenTime string) (time.Duration, error) {
	thenDate, err := thenDate(thenDay, thenTime)
	if err != nil {
		return 0, err
	}
	return thenDate.Sub(time.Now()), nil
}

func thenDate(thenDay, thenTime string) (time.Time, error) {
	cmd := exec.Command("date", "--rfc-2822", "--date", thenDay+" "+thenTime)
	var cmdResult bytes.Buffer
	cmd.Stdout = &cmdResult
	err := cmd.Run()
	if err != nil {
		return time.Now(), err
	}
	then, err := time.Parse(time.RFC1123Z, strings.TrimSuffix(cmdResult.String(), "\n"))
	if err != nil {
		return time.Now(), err
	}
	return then, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, data *templateData) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		t, _ = template.ParseFiles(data.DocumentRoot + "/" + data.Path + "/" + tmpl + ".html")
	}
	w.WriteHeader(http.StatusOK)
	t.Execute(w, data)
}
