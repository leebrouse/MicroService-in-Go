package main

import (
	"net/http"

	"github.com/leebrouse/MicroService-in-Go/log-serivce/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// insert logs into mongo and sent respose to the broker-service
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	//TODO:
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	//insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
	}

	//return response to the broker
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
