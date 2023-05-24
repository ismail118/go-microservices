package main

import (
	"fmt"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var reqPayload mailMessage

	err := app.readJson(w, r, &reqPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	msg := Message{
		From:    reqPayload.From,
		To:      reqPayload.To,
		Subject: reqPayload.Subject,
		Data:    reqPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("sent to %s", reqPayload.To),
	}

	err = app.writeJson(w, http.StatusAccepted, payload)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
}
