package main

import (
	"net/http"
)

func AddReminder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		thenDay := r.URL.Query().Get("day")
		if thenDay == "" {
			http.Error(w, "missing day", http.StatusBadRequest)
			return
		}
		thenTime := r.URL.Query().Get("time")
		if thenTime == "" {
			http.Error(w, "missing time", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(thenDay))
	})
}

func nowDate() { return }

func thenDate() { return }
