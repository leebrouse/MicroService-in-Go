package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// modify the webPort from 80 to 8080
const webPort = "8080"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	//connect to the rabbit
	conn, err := connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	app := Config{
		Rabbit: conn,
	}

	log.Println("Starting broker service on port: ", webPort)

	//Config server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// Start server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// function to connect the rabbitMQ
func connect() (*amqp.Connection, error) {

	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		// try to connect
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = conn
			break
		}

		//try 5 times ,if can't connect with the rabbitMQ
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off ....")
		time.Sleep(backOff)
		continue
	}

	return connection, nil

}
