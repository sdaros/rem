package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Reminder struct {
	*App
	ReminderMsg string
	SuccessMsg  string
	ThenDate    string
	ThenTime    string
	InputType   string
}

func CreateReminder(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reminder := newReminder(app)
		switch r.Method {
		case "GET":
			renderTemplate(w, r, reminder)
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
		InputType:  "time",
	}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, rem *Reminder) {
	templateFile := rem.DocumentRoot + "/" + rem.Path + "/create.html"
	t, err := template.ParseFiles(templateFile)
	die("error when rendering template: %v", err)
	if matchesAndroidBrowserUserAgent(r) {
		fallbackToFormInputTypeText(rem)
	}
	w.WriteHeader(http.StatusOK)
	t.Execute(w, rem)
}

func submitReminder(w http.ResponseWriter, r *http.Request, rem *Reminder) {
	rem.ThenDate = r.FormValue("date")
	rem.ThenTime = r.FormValue("time")
	rem.ReminderMsg = r.FormValue("message")
	delay, err := timeToSleepFor(rem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendReminderAfter(delay, rem)
	logNewReminder(delay, rem)
	if isAjaxRequest(r) {
		w.WriteHeader(http.StatusOK)
		return
	}
	renderTemplate(w, r, rem)
	return
}

func sendReminderAfter(delay time.Duration, rem *Reminder) {
	go func(time.Duration, *Reminder) {
		select {
		case <-time.After(delay):
			sendReminder(rem)
		}
	}(delay, rem)
	return
}

func logNewReminder(delay time.Duration, rem *Reminder) {
	thenDateTime := rem.ThenDate + " " + rem.ThenTime
	rem.SuccessMsg = "Thank you! Your reminder will be sent at " + thenDateTime
	log.Printf("The reminder '%v' will be sent at %v", rem.ReminderMsg, thenDateTime)
	return
}

func timeToSleepFor(rem *Reminder) (time.Duration, error) {
	thenAsUtc, err := time.Parse("2006-01-02 15:04", rem.ThenDate+" "+rem.ThenTime)
	if err != nil {
		return 0, err
	}
	nowAsUtc, _ := time.Parse(time.ANSIC, time.Now().Format(time.ANSIC))
	return thenAsUtc.Sub(nowAsUtc), nil
}

func sendReminder(rem *Reminder) {
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

func isAjaxRequest(r *http.Request) bool {
	if r.Header.Get("X-Requested-With") == "XMLHTTPRequest" {
		return true
	}
	return false
}

func matchesAndroidBrowserUserAgent(r *http.Request) bool {
	if strings.Contains(r.Header.Get("User-Agent"), "Android") {
		return true
	}
	return false
}

func fallbackToFormInputTypeText(rem *Reminder) {
	// Some versions of android's default browser do
	// not handle <input type="time"> properly.
	rem.InputType = "text"
}
