package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/leebrouse/MicroService-in-Go/log-serivce/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Port
const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"

	/* URL Just for Testing */
	// mongoURL = "mongodb://localhost:27017"

	grpcPort = "50001"
)

// mongo-client
var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	//create a context in order to disconnect(within 15 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//call disconnect
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	//create app config
	app := Config{
		Models: data.New(client),
	}

	/*  Blocking function:

	var forever chan struct{}
	<-forever

	*/

	var forever chan struct{}
	//start web server
	go app.serve()
	<-forever

}

func (app *Config) serve() {

	log.Printf("Starting log-service in:%s \n", webPort)

	//set server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.router(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection option
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:" + err.Error())
		return nil, err
	}

	log.Println("connect")

	return conn, nil
}
