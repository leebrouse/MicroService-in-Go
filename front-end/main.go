package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 自定义HTML渲染
	// r.SetHTMLTemplate(loadTemplates())
	r.LoadHTMLGlob("./cmd/web/templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.page.html", nil)
	})

	fmt.Println("Starting front end service on port 80")
	if err := r.Run(":80"); err != nil {
		log.Panic(err)
	}
}

// 加载模板
// func loadTemplates() *template.Template {
// 	partials := []string{
// 		"./cmd/web/templates/base.layout.html",
// 		"./cmd/web/templates/header.partial.html",
// 		"./cmd/web/templates/footer.partial.html",
// 	}

// 	// 需要解析的页面
// 	pages := []string{
// 		"./cmd/web/templates/test.page.html",
// 	}

// 	// 合并所有模板
// 	allTemplates := append(pages, partials...)

// 	tmpl, err := template.ParseFiles(allTemplates...)
// 	if err != nil {
// 		log.Panicf("Failed to load templates: %v", err)
// 	}
// 	return tmpl
// }
