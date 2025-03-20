package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) router() http.Handler {
	//create a new router
	mux := chi.NewRouter()

	//set router limitors
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	//add a middleware that can be a load balancers
	mux.Use(middleware.Heartbeat("/ping"))

	//add log-service handler
	mux.Post("/log", app.WriteLog)

	return mux
}
