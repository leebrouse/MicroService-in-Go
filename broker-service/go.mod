module github.com/leebrouse/MicroService-in-Go/broker-service

go 1.24.1

// add 3 third-package
// 1.go get github.com/go-chi/chi/v5
// 2.go get github.com/go-chi/chi/v5/middleware
// 3.go get github.com/go-chi/cors

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.1
	github.com/golang/protobuf v1.5.4
	github.com/rabbitmq/amqp091-go v1.10.0
	google.golang.org/grpc v1.71.0
)

require (
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/protobuf v1.36.4 // indirect
)
