package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", func(c *gin.Context) {
			name := c.DefaultPostForm("name", "guest")
			c.String(http.StatusOK, "Hello "+name)
		})
		v1.POST("/submit", func(c *gin.Context) {
			name := c.DefaultPostForm("name", "guest")
			c.String(http.StatusOK, "Hello "+name)
		})
		v1.POST("/read", func(c *gin.Context) {
			name := c.DefaultPostForm("name", "guest")
			c.String(http.StatusOK, "Hello "+name)
		})
	}

	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", func(c *gin.Context) {
			name := c.DefaultPostForm("name", "guest")
			c.String(http.StatusOK, "Hello "+name)
		})
		v2.POST("/submit", func(c *gin.Context) {
			name := c.DefaultPostForm("name", "guest")
			c.String(http.StatusOK, "Hello "+name)
		})
		v2.POST("/read", func(c *gin.Context) {
			name := c.DefaultPostForm("name", "guest")
			c.String(http.StatusOK, "Hello "+name)
		})
	}

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Error while running router: %v", err)
	}
}
