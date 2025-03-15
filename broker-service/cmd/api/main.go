package main

import "github.com/gin-gonic/gin"

const webPort = ":80"

type Config struct{}

func NewEngin() *gin.Engine {
	r := gin.Default()

	//add routers middleware and so on
	Routers(r)

	return r
}

func main() {
	r := NewEngin()

	r.Run(webPort)
}
