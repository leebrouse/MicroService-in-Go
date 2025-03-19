package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResonse{
		Error:   false,
		Message: "Hit the broker",
	}

	app.writeJSON(w, http.StatusOK, payload)
}
