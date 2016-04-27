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

func CreateReminder(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isQueryParamsEmpty(r) {
			// then return a form to create a new reminder
			w.WriteHeader(http.StatusOK)
			renderTemplate(w, "create", nil)
			return
		}
		submitReminder(w, r, app)
		return
	})
}

func CreateReminderViaForm(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		submitReminder(w, r, app)
		type data struct {
			ReminderSuccess string
		}
		renderTemplate(w, "create",
			&data{"Your reminder has been added, thank you!"})
		return
	})
}

func submitReminder(w http.ResponseWriter, r *http.Request, app *App) {
	thenDay := r.URL.Query().Get("day")
	if thenDay == "" {
		thenDay = r.FormValue("day")
	}
	thenTime := r.URL.Query().Get("time")
	if thenTime == "" {
		thenTime = r.FormValue("time")
	}
	message := r.URL.Query().Get("message")
	if message == "" {
		message = r.FormValue("message")
	}
	delay, err := timeToSleepFor(thenDay, thenTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("The reminder '%v' will be sent to you at %v", message, time.Now().Add(delay))
	go func(time.Duration, *App, string) {
		select {
		case <-time.After(delay):
			execute(app, message)
		}
	}(delay, app, message)
}

func execute(app *App, msg string) {
	cmd := exec.Command(app.Lookup("home").(string)+"/bin/rem_script", msg)
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

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, data)
}

func isQueryParamsEmpty(r *http.Request) bool {
	thenTime := r.URL.Query().Get("time")
	message := r.URL.Query().Get("message")
	if thenTime == "" || message == "" {
		return true
	}
	return false
}
