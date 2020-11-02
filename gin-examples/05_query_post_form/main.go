package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// POST /post?id=1234&page=1 HTTP/1.1
	// Content-Type: application/x-www-form-urlencoded
	// name=manu&message=this_is_great
	router.POST("/post", func(c *gin.Context) {

		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")

		c.String(http.StatusOK, fmt.Sprintf("id: %s; page: %s; name: %s; message: %s", id, page, name, message))
	})

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Error while running router: %v", err)
	}
}
