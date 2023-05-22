package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var reqPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &reqPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(reqPayload.Email)
	if err != nil {
		app.errorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(reqPayload.Password)
	if err != nil || !valid {
		app.errorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	err = app.writeJson(w, http.StatusAccepted, payload)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
}
