package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/leebrouse/MicroService-in-Go/log-serivce/data"
	"github.com/leebrouse/MicroService-in-Go/log-serivce/logs"
	"google.golang.org/grpc"
)

// UnimplementedLogServiceServer must be embedded to have
// forward compatible implementations.
type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, in *logs.LogRequest) (*logs.LogResponse, error) {
	// get logEntry
	input := in.GetLogEntry()

	//write log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// insert data
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		return &logs.LogResponse{Result: "Write Log Failed"}, err
	}

	// return response
	return &logs.LogResponse{Result: "Write Log via gRpc successed"}, nil
}

// gRpc Server listener
func (app *Config) gRpcListen() {
	// create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// create grpc server
	grpcServer := grpc.NewServer()

	// register service
	logs.RegisterLogServiceServer(grpcServer, &LogServer{Models: app.Models})

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
