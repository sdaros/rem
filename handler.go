package main

import (
	"html/template"
	"log"
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

func (self *Reminder) sendReminder(r *http.Request) {
	delay, err := self.calculateNotificationDelay(r)
	if err != nil {
		self.notifyAndLogTheError("error: unable to calculate delay:  " + err.Error())
		return
	}
	self.logNewReminder()
	select {
	case <-time.After(delay):
		err := self.Notification.Notify()
		die("error: unable to use Notification API: %v", err)
	}
	return
}

func (self *Reminder) calculateNotificationDelay(r *http.Request) (time.Duration, error) {
	clientNow, err := time.Parse("Mon Jan _2 2006 15:04:05 GMT-0700 (MST)",
		r.FormValue("client-now"))
	if err != nil {
		return 0, err
	}
	then, err := time.ParseInLocation("2006-01-02 15:04",
		self.ThenDate+" "+self.ThenTime, clientNow.Location())
	if err != nil {
		return 0, err
	}
	now := time.Now().In(clientNow.Location())
	delay := then.Sub(now)
	return delay, nil
}

func (self *Reminder) renderTemplate(w http.ResponseWriter) {
	tmpl, err := template.New("createReminder").Parse(createReminderTemplate)
	die("error: unable to render template: %v", err)
	w.WriteHeader(http.StatusOK)
	err = tmpl.ExecuteTemplate(w, "createReminder", self.TemplateData)
	die("error: unable to execute template: %v", err)
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

func (self *Reminder) notifyAndLogTheError(err string) {
	log.Printf("%v", err)
	self.Notification.Message = err
	notifyErr := self.Notification.Notify()
	die("%v", notifyErr)
}

func (self *Reminder) logNewReminder() {
	log.Printf("The reminder '%v' will be sent at %v",
		self.Message, self.ThenDate+" "+self.ThenTime)
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
