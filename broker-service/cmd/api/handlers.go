package main

import (
	"log"
	"net/http"
)

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
