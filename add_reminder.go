package main

import (
	"bytes"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func AddReminder(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isQueryParamsEmpty(r) {
			path := app.Lookup("path").(string)
			http.Redirect(w, r, path+"/create/", http.StatusFound)
		}
		thenDay := r.URL.Query().Get("day")
		thenTime := r.URL.Query().Get("time")
		message := r.URL.Query().Get("message")
		delay, err := timeToSleepFor(thenDay, thenTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Your reminder has been added. Thank you."))
		log.Printf("The reminder '%v' will be sent to you at %v", message, time.Now().Add(delay))
		go func(time.Duration, *App, string) {
			select {
			case <-time.After(delay):
				execute(app, message)
			}
		}(delay, app, message)
	})
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

func isQueryParamsEmpty(r *http.Request) bool {
	thenTime := r.URL.Query().Get("time")
	message := r.URL.Query().Get("message")
	if thenTime == "" || message == "" {
		return true
	}
	return false
}
