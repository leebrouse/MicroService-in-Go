package main

import (
	"bytes"
	"encoding/json"
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

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
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

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string
		Data string
	}
	entry.Name = name
	entry.Data = data

	//create some json we`ll sent ro the auth microservice
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	//call the service
	request, err := http.NewRequest("POST", "http://log-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	//Create a new client to POST the request
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
