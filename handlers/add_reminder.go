package handlers

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func AddReminder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		thenDay := r.URL.Query().Get("day")
		thenTime := r.URL.Query().Get("time")
		message := r.URL.Query().Get("message")
		if thenTime == "" {
			http.Error(w, "Example: provide a `time=1330` param", http.StatusBadRequest)
			return
		}
		if message == "" {
			http.Error(w, "Example: provide a `message=buy-beer` param", http.StatusBadRequest)
			return
		}
		delay, err := timeToSleepFor(thenDay, thenTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Your reminder has been added. Thank you."))
		go func(time.Duration, string) {
			select {
			case <-time.After(delay):
				execute(message)
			}
		}(delay, message)
	})
}

func timeToSleepFor(thenDay, thenTime string) (time.Duration, error) {
	then, err := thenAsUnix(thenDay, thenTime)
	if err != nil {
		return 0, err
	}
	now, err := nowAsUnix()
	if err != nil {
		return 0, err
	}
	return (time.Duration(then-now) * time.Second), nil
}

func nowAsUnix() (int64, error) {
	cmd := exec.Command("date", "+%s")
	var cmdResult bytes.Buffer
	cmd.Stdout = &cmdResult
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	now, err := cmdResultToInt(cmdResult)
	if err != nil {
		return 0, err
	}
	return now, nil
}

func thenAsUnix(thenDay, thenTime string) (int64, error) {
	cmd := exec.Command("date", "--date", thenDay+" "+thenTime, "+%s")
	var cmdResult bytes.Buffer
	cmd.Stdout = &cmdResult
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	then, err := cmdResultToInt(cmdResult)
	if err != nil {
		return 0, err
	}
	return then, nil
}

func cmdResultToInt(cmdResult bytes.Buffer) (int64, error) {
	result, err := cmdResult.ReadString('\n')
	if err != nil {
		return 0, err
	}
	resultAsUnixTime, err := strconv.Atoi(strings.TrimSuffix(result, "\n"))
	if err != nil {
		return 0, err
	}
	return int64(resultAsUnixTime), nil

}

func execute(msg string) {
	cmd := exec.Command(os.Getenv("HOME")+"/bin/rem_script", msg)
	var cmdResult bytes.Buffer
	cmd.Stdout = &cmdResult
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return
}
