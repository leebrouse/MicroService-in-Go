package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/leebrouse/MicroService-in-Go/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect rabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	//start listening for message
	log.Println("Listening for and consuming RabbitMQ messages...")

	//create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		fmt.Println("Can't create a  new consumer")
		return
	}

	//watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	
}

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
