package main

import (
	"errors"
	"net/http"
	"net/url"
)

type (
	Notification struct {
		Api     string `json: "NotificationApi"`
		Token   string `json: "ApiToken"`
		User    string `json: "ApiUser"`
		Message string
	}
)

func (self *Notification) Notify() error {
	resp, err := http.PostForm(self.Api,
		url.Values{"token": {self.Token},
			"user":    {self.User},
			"message": {self.Message}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		err := "error when using Notification API, got: " + resp.Status
		return errors.New(err)
	}
	return nil
}
