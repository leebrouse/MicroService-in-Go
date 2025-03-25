package main

import (
	"context"
	"log"
	"time"

	"github.com/leebrouse/MicroService-in-Go/log-serivce/data"
)

type RpcServer struct{}

type RpcPayload struct {
	Name string
	Data string
}

// Handler to operate the log function via rpc
func (r *RpcServer) LogInfo(payload RpcPayload, resp *string) error {
	// get collection in mongo
	collection := client.Database("logs").Collection("logs")
	// collect insert
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:     payload.Name,
		Data:     payload.Data,
		CreateAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing messge to mongo", err)
		return err
	}

	*resp = "Process payload via RPC:" + payload.Name
	return nil
}


