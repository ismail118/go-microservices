package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	default:
		app.errorJson(w, errors.New("unknown action"))
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
