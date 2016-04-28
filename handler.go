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
	DefaultTime     string
	ReminderSuccess string
	Path            string
}

func CreateReminder(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := initialiseTemplateData(w, app)
		switch r.Method {
		case "GET":
			renderTemplate(w, "create", data, app)
			return
		case "POST":
			submitReminder(w, r, data, app)
			return
		}
	})
}

func initialiseTemplateData(w http.ResponseWriter, app *App) *templateData {
	config := app.Lookup("config").(*Config)
	reminderSuccess := ""
	// TODO: +30m
	thirtyMinutesFromNow := time.Now().Add(time.Duration(30) * time.Minute)
	defaultTime := thirtyMinutesFromNow.Format("15:04")
	return &templateData{defaultTime, reminderSuccess, config.Path}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}, app *App) {
	config := app.Lookup("config").(*Config)
	t, _ := template.ParseFiles(config.ReminderTemplate)
	w.WriteHeader(http.StatusOK)
	t.Execute(w, data)
}

func submitReminder(w http.ResponseWriter, r *http.Request, data *templateData, app *App) {
	thenDay := r.FormValue("day")
	thenTime := r.FormValue("time")
	message := r.FormValue("message")
	delay, err := timeToSleepFor(thenDay, thenTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data.ReminderSuccess = "Your reminder has been added, thank you!"
	renderTemplate(w, "create", data, app)
	log.Printf("The reminder '%v' will be sent to you at %v",
		message, time.Now().Add(delay).Format(time.RFC3339))
	go func(time.Duration, *App, string) {
		select {
		case <-time.After(delay):
			execute(app, message)
		}
	}(delay, app, message)
}

func timeToSleepFor(thenDay, thenTime string) (time.Duration, error) {
	thenDate, err := thenDate(thenDay, thenTime)
	if err != nil {
		return 0, err
	}
	return thenDate.Sub(time.Now()), nil
}

func thenDate(thenDay, thenTime string) (time.Time, error) {
	cmd := exec.Command("date", "--rfc-3339=seconds", "--date", thenDay+" "+thenTime)
	var cmdResult bytes.Buffer
	cmd.Stdout = &cmdResult
	err := cmd.Run()
	if err != nil {
		return time.Now(), err
	}
	then, err := time.Parse("2006-01-02 15:04:05-07:00",
		strings.TrimSuffix(cmdResult.String(), "\n"))
	if err != nil {
		return time.Now(), err
	}
	return then, nil
}

func execute(app *App, msg string) {
	config := app.Lookup("config").(*Config)
	cmd := exec.Command(config.RemScript, msg)
	var cmdResult bytes.Buffer
	cmd.Stdout = &cmdResult
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return
}
