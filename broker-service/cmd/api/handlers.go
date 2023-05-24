package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Success",
	}

	err := app.writeJson(w, http.StatusOK, payload)
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var reqPayload RequestPayload

	err := app.readJson(w, r, &reqPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	switch reqPayload.Action {
	case "auth":
		app.authenticate(w, reqPayload.Auth)
	case "log":
		app.logItem(w, reqPayload.Log)
	case "mail":
		app.sendMail(w, reqPayload.Mail)
	default:
		app.errorJson(w, errors.New("unknown action"))
	}

}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	logServiceURL := "http://logger-service/log"

	req, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling logger service"), res.StatusCode)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logger"

	err = app.writeJson(w, http.StatusAccepted, payload)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll to the auth microservice
	jsonData, err := json.Marshal(a)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	// call the service
	req, err := http.NewRequest(http.MethodPost, "http://authentication-service/authenticate",
		bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer res.Body.Close()

	// make sure we get back the correct status code
	if res.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentials"))
		return
	} else if res.StatusCode == http.StatusBadRequest {
		app.errorJson(w, errors.New("bad request"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling auth service"))
		return
	}

	// read response.body
	var jsonFromService jsonResponse

	err = json.NewDecoder(res.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	if jsonFromService.Error {
		app.errorJson(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	err = app.writeJson(w, http.StatusAccepted, payload)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	mailServiceURL := "http://mail-service/send"
	req, err := http.NewRequest(http.MethodPost, mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling mail service"), res.StatusCode)
		return
	}

	resPayload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("mail send to %s", msg.To),
	}

	err = app.writeJson(w, http.StatusAccepted, resPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
}
