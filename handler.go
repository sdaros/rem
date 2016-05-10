package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type Reminder struct {
	*Notification
	TemplateData *Template
	ThenDate     string
	ThenTime     string
}

func CreateReminder(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reminder := newReminder(app)
		switch r.Method {
		case "GET":
			if matchesAndroidBrowserUserAgent(r) {
				reminder.TemplateData.fallbackToFormInputTypeText()
			}
			reminder.renderTemplate(w)
			return
		case "POST":
			reminder.submit(w, r)
			return
		}
	})
}

func newReminder(app *App) *Reminder {
	return &Reminder{
		Notification: app.Notification,
		TemplateData: &Template{
			SuccessMsg: "",
			InputType:  "time",
			Domain:     app.Domain,
			Path:       app.Path,
		},
	}
}

func (self *Reminder) renderTemplate(w http.ResponseWriter) {
	tmpl, err := template.New("createReminder").Parse(createReminderTemplate)
	die("error when rendering template: %v", err)
	w.WriteHeader(http.StatusOK)
	err = tmpl.ExecuteTemplate(w, "createReminder", self.TemplateData)
	die("error when executing template: %v", err)
	return
}

func (self *Reminder) submit(w http.ResponseWriter, r *http.Request) {
	if err := self.validateClientInput(r); err != nil {
		http.Error(w, "Sorry, we were unable to process your Input.",
			http.StatusInternalServerError)
		return
	}
	go self.sendReminder(r)
	if isAjaxRequest(r) {
		w.WriteHeader(http.StatusOK)
		return
	}
	self.TemplateData.SuccessMsg = "Thank you! Your reminder will be sent at " +
		self.ThenDate + " " + self.ThenTime
	self.renderTemplate(w)
	return
}

func (self *Reminder) validateClientInput(r *http.Request) error {
	self.ThenDate = r.FormValue("date")
	self.ThenTime = r.FormValue("time")
	self.Message = r.FormValue("message")
	_, err := time.Parse("2006-01-02 15:04", self.ThenDate+" "+self.ThenTime)
	if err != nil {
		return err
	}
	return nil
}

func (self *Reminder) sendReminder(r *http.Request) {
	loc := self.detectClientLocation(r)
	delay := self.calculateNotificationDelay(loc)
	select {
	case <-time.After(delay):
		err := self.Notification.Notify()
		die("error when using Notification API, got: %v", err)
	}
	return
}

func (self *Reminder) detectClientLocation(r *http.Request) (location string) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	resp, err := http.Get("http://freegeoip.net/json/" + ip)
	if err != nil {
		self.notifyTheError("error when trying to geolocate client IP: " +
			err.Error())
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	type freegeoipJson struct {
		Time_zone string
	}
	jsonResp := new(freegeoipJson)
	if err = decoder.Decode(&jsonResp); err != nil {
		self.notifyTheError("error when trying to parse json response: " + err.Error())
	}
	if jsonResp.Time_zone == "" {
		self.notifyTheError("error: timezone from freegeoip was empty!" + err.Error())
	}
	return jsonResp.Time_zone
}

func (self *Reminder) calculateNotificationDelay(loc string) time.Duration {
	locationOfClient, err := time.LoadLocation(loc)
	if err != nil {
		self.notifyTheError("error when using autodetected location: " + err.Error())
	}
	then, err := time.ParseInLocation("2006-01-02 15:04",
		self.ThenDate+" "+self.ThenTime, locationOfClient)
	if err != nil {
		self.notifyTheError("error when parsing autodetected location: " + err.Error())
		return 0
	}
	now := time.Now().In(locationOfClient)
	delay := then.Sub(now)
	logReminder(now, delay, self.Message)
	return delay
}

func logReminder(now time.Time, delay time.Duration, msg string) {
	log.Printf("The reminder '%v' will be sent at %v", msg, now.Add(delay))
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

func (self *Reminder) notifyTheError(err string) {
	self.Notification.Message = err
	notifyErr := self.Notification.Notify()
	die("%v", notifyErr)
}
