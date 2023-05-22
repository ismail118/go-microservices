package main

import (
	"log-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var reqPayload JSONPayload

	err := app.readJson(w, r, &reqPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	// insert data
	event := data.LogEntry{
		Name: reqPayload.Name,
		Data: reqPayload.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	err = app.writeJson(w, http.StatusAccepted, resp)
	if err != nil {
		app.errorJson(w, err)
		return
	}
}
