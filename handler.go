package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net"
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
	if err := validateClientInput(r, rem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go sendReminder(r, rem)
	if isAjaxRequest(r) {
		w.WriteHeader(http.StatusOK)
		return
	}
	rem.SuccessMsg = "Thank you! Your reminder will be sent at " +
		rem.ThenDate + " " + rem.ThenTime
	renderTemplate(w, r, rem)
	return
}

func validateClientInput(r *http.Request, rem *Reminder) error {
	rem.ThenDate = r.FormValue("date")
	rem.ThenTime = r.FormValue("time")
	rem.ReminderMsg = r.FormValue("message")
	_, err := time.Parse("2006-01-02 15:04", rem.ThenDate+" "+rem.ThenTime)
	if err != nil {
		return err
	}
	return nil
}

func sendReminder(r *http.Request, rem *Reminder) {
	loc := make(chan string)
	loc = detectClientLocation(r)
	log.Printf("loc: %v", loc)
	delay, err := calculateDelay(loc, rem)
	die("error parsing datetime provided by client: %v", err)
	go func(time.Duration, *Reminder) {
		select {
		case <-time.After(delay):
			sendReminderToNotificationApi(rem)
		}
	}(delay, rem)
	return
}

func detectClientLocation(r *http.Request) (location chan string) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	resp, err := http.Get("https://freegeoip.net/json/" + ip)
	die("error while trying to geolocate client IP: %v", err)
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	type freegeoipJson struct {
		time_zone string
	}
	jsonResp := new(freegeoipJson)
	err = decoder.Decode(&jsonResp)
	die("error while trying to parse json response: %v", err)
	location <- jsonResp.time_zone
	return
}

func calculateDelay(loc <-chan string, rem *Reminder) (time.Duration, error) {
	locationOfClient, err := time.LoadLocation(<-loc)
	die("error while parsing freegeoip location: %v", err)
	then, err := time.ParseInLocation("2006-01-02 15:04",
		rem.ThenDate+" "+rem.ThenTime, locationOfClient)
	if err != nil {
		return 0, err
	}
	now := time.Now().In(locationOfClient)
	delay := then.Sub(now)
	logNewReminder(delay, rem)
	return delay, nil
}

func sendReminderToNotificationApi(rem *Reminder) {
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

func logNewReminder(delay time.Duration, rem *Reminder) {
	thenDateTime := rem.ThenDate + " " + rem.ThenTime
	log.Printf("The reminder '%v' will be sent at %v", rem.ReminderMsg, thenDateTime)
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
