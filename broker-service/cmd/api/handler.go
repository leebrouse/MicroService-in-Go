package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// Main request struct
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth"`
}

// Sub struct
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Broker-Service handler
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// Call Auth-Service handler
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	//read request
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authentication(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("Unknow Action"))
	}
}

func (app *Config) authentication(w http.ResponseWriter, a AuthPayload) {
	//create some json we`ll sent ro the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authentication", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//Create a new client to POST the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	//make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Error calling auth service"))
		return
	}

	//create a varable we'll read response.body into
	var jsonFromService jsonResponse

	//decode the json from the auth service and input to the jsonFromService struct
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//Check the error message in the jsonFromService body
	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	//All the operation are done ,now output the response
	payload := jsonResponse{
		Error:   false,
		Message: "Authentication",
		Data:    jsonFromService.Data,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}
