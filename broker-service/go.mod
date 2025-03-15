module github.com/leebrouse/MicroService-in-Go/broker-service

go 1.24.1

// add 3 third-package
// 1.go get github.com/go-chi/chi/v5
// 2.go get github.com/go-chi/chi/v5/middleware
// 3.go get github.com/go-chi/cors

require (
	github.com/go-chi/chi/v5 v5.2.1 // indirect
	github.com/go-chi/cors v1.2.1 // indirect
)
