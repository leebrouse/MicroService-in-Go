package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	//web templates
	r.LoadHTMLGlob("./cmd/web/templates/*")

	//default router
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.page.html", nil)
	})

	return r
}

func main() {
	r := NewRouter()

	fmt.Println("Starting front end service on port 80")
	if err := r.Run(":80"); err != nil {
		log.Panic(err)
	}
}
