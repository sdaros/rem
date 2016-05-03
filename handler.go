package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

type Reminder struct {
	*App
	DefaultTime string
	ReminderMsg string
	SuccessMsg  string
	ThenDay     string
	ThenTime    string
}

func CreateReminder(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reminder := newReminder(app)
		reminder.DefaultTime = DefaultTime(r)
		switch r.Method {
		case "GET":
			renderTemplate(w, reminder)
			return
		case "POST":
			submitReminder(w, r, reminder)
			return
		}
	})
}

func newReminder(app *App) *Reminder {
	return &Reminder{
		App:        app,
		SuccessMsg: "",
	}
}

func renderTemplate(w http.ResponseWriter, rem *Reminder) {
	templateFile := rem.DocumentRoot + "/" + rem.Path + "/create.html"
	t, err := template.ParseFiles(templateFile)
	die("error when rendering template: %v", err)
	w.WriteHeader(http.StatusOK)
	t.Execute(w, rem)
}

func submitReminder(w http.ResponseWriter, r *http.Request, rem *Reminder) {
	rem.ThenDay = r.FormValue("day")
	rem.ThenTime = r.FormValue("time")
	rem.ReminderMsg = r.FormValue("message")
	delay, err := timeToSleepFor(rem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func(time.Duration, *Reminder) {
		select {
		case <-time.After(delay):
			execute(rem)
		}
	}(delay, rem)

	thenDateTime := time.Now().Add(delay).Format("2006-01-02 15:04:05-07:00")
	rem.SuccessMsg = "Thank you! Your reminder will be sent at " + thenDateTime
	renderTemplate(w, rem)
	log.Printf("The reminder '%v' will be sent at %v", rem.ReminderMsg, thenDateTime)
}

func timeToSleepFor(rem *Reminder) (time.Duration, error) {
	thenDate, err := thenDate(rem)
	if err != nil {
		return 0, err
	}
	return thenDate.Sub(time.Now()), nil
}

func thenDate(rem *Reminder) (time.Time, error) {
	cmd := exec.Command("date", "--rfc-3339=seconds", "--date",
		rem.ThenDay+" "+rem.ThenTime)
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

func execute(rem *Reminder) {
	resp, err := http.PostForm(rem.NotificationApi,
		url.Values{"token": {rem.ApiToken},
			"user":    {rem.ApiUser},
			"message": {rem.ReminderMsg}})
	die("error when using Notification API: %v", err)
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		die("error when using Notification API: %v",
			errors.New(resp.Status))
	}
	return
}
