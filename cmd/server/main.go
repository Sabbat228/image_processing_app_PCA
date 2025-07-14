package main

import (
	"image-processing-app/internal/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	if err := os.MkdirAll("web/templates", os.ModePerm); err != nil {
		log.Fatal("Failed to create templates directory:", err)
	}
	if err := os.MkdirAll("web/static/uploads", os.ModePerm); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}

	handlers.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
