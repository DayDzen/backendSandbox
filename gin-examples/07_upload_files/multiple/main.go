package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			dst := fmt.Sprintf("upload_files/multiple/files/%v", file.Filename)
			err := c.SaveUploadedFile(file, dst)
			if err != nil {
				log.Fatalf("Error while saving uploaded file: %v", err)
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Error while running router: %v", err)
	}
}
