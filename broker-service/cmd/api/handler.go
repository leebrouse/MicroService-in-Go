package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BrokerInterface interface {
	BrokerService(c *gin.Context)
}

type jsonResonse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Broker struct{}

// broker handler

func NewBrokerHandler() BrokerInterface {
	return &Broker{}
}

func (b *Broker) BrokerService(c *gin.Context) {

	payload := &jsonResonse{
		Error:   false,
		Message: "Hit the broker",
	}

	c.JSON(http.StatusOK, payload)
}
