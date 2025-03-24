package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// create new cosumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
		// queueName: ,
	}

	//set up Exchange in the consumer
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup channel
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// request from the rabbitmq
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	//set up queue
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
	}

	message, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// continous listening
	var forever chan struct{}
	go func() {
		for d := range message {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()
	fmt.Printf("Waiting for message [Exchange,Queue] [logs_topic,%s]\n", q.Name)
	<-forever

	return nil
}

// Function handlepayload for distrubting the payload request to the specific server
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
			return
		}
	case "auth":
		// auth service
		err := authEvent(payload)
		if err != nil {
			log.Println(err)
			return
		}
	default:
		// call logEvent
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

// Remot call auth-service
func authEvent(entry Payload) error {
	//create some json we`ll sent ro the auth microservice
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	//call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authentication", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	//Create a new client to POST the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		return err
	} else if response.StatusCode != http.StatusAccepted {
		return err
	}

	//create a varable we'll read response.body into
	type jsonResponse struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Data    any    `json:"data,omitempty"`
	}

	var jsonFromService jsonResponse

	//decode the json from the auth service and input to the jsonFromService struct
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		return err
	}

	//Check the error message in the jsonFromService body
	if jsonFromService.Error {
		return err
	}

	return nil
}

// Remote call log-service
func logEvent(entry Payload) error {
	//create some json we`ll sent ro the auth microservice
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	//call the service
	request, err := http.NewRequest("POST", "http://log-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	//Create a new client to POST the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
