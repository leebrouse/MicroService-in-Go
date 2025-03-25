package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"github.com/leebrouse/MicroService-in-Go/broker-service/event"
	"github.com/leebrouse/MicroService-in-Go/broker-service/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Main request struct
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

// Sub struct
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
	case "log":
		// app.logItem(w,requestPayload.Log)
		// app.logEventViaRabbit(w, requestPayload.Log)
		app.logEventViaRpc(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("Unknow Action"))
	}
}

// Remote call auth-service by REST API
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

// Remote call log-service by REST API
func (app *Config) logItem(w http.ResponseWriter, l LogPayload) {
	//create some json we`ll sent ro the auth microservice
	jsonData, _ := json.MarshalIndent(l, "", "\t")

	//call the service
	request, err := http.NewRequest("POST", "http://log-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	//Create a new client to POST the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	//make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Error calling log service"))
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
		Message: "logged",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

// Remote call mail-service by REST API
func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {
	// TODO:
	//create some json we`ll sent ro the auth microservice
	jsonData, _ := json.MarshalIndent(m, "", "\t")

	//call the service
	request, err := http.NewRequest("POST", "http://mail-service/send", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	//Create a new client to POST the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	//make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Error calling mail service"))
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
		Message: "Message sent to " + m.To,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// push message to the rabbitmq queue
func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)
}

// push message to the rabbit queue
func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(payload, "", "\t")

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

// RpcPayload struct
type RpcPayload struct {
	Name string
	Data string
}

// Remote call RpcServer via RPC
func (app *Config) logEventViaRpc(w http.ResponseWriter, l LogPayload) {
	// rpc dial
	client, err := rpc.Dial("tcp", "log-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayload := RpcPayload{
		Name: l.Name,
		Data: l.Data,
	}
	var result string

	//rpc call the logInfo function
	client.Call("RpcServer.LogInfo", rpcPayload, &result)

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

// Remote Call RpcServer via gRPC
func (app *Config) logEventViaGrpc(w http.ResponseWriter, r *http.Request) {
	var logPayload LogPayload
	// read request
	err := app.readJSON(w, r, &logPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// create grpc client
	conn, err := grpc.NewClient("log-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	// call grpc server
	client := logs.NewLogServiceClient(conn)

	// create ctx to limit the operating time
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// write log call to the grpc server
	_, err = client.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: logPayload.Name,
			Data: logPayload.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// output the result
	payload := jsonResponse{
		Error:   false,
		Message: "Logged via gRPC",
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}
