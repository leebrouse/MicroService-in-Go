package main

import (
	"fmt"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	//decode the request
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// create message
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	// call send mail function
	err = app.Mail.SendSMTPMessage(msg)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//output the payload
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintln("Sent to ", msg.To),
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}
