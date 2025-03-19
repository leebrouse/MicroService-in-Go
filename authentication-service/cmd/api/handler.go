package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authentication(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
	}

	/* Validate the user against the database */
	//1. Geting the username by the email
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	//2. Matching the user password with the request password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	//3. Create Payload struct jsonResponse
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s \n", user.Email),
		Data:    user,
	}

	//4. Writing the jsonResponse and send the 202 as statusCode
	app.writeJSON(w, http.StatusAccepted, payload)
}
