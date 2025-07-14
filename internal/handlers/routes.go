package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image-processing-app/internal/services"
	"log"
	"net/http"
)

func SetupRoutes(r *gin.Engine) {
	imageProcessor := services.NewImageProcessor(50, 200)
	uploadsDir := "web/static/uploads"

	homeHandler := NewHomeHandler()
	uploadHandler := NewUploadHandler(uploadsDir)
	resultsHandler := NewResultsHandler(uploadsDir)
	analyzeHandler := NewAnalyzeHandler(imageProcessor, uploadsDir)

	r.GET("/", homeHandler.Handle)
	r.POST("/upload", uploadHandler.Handle)
	r.GET("/results/:id", resultsHandler.Handle)
	r.POST("/analyze", analyzeHandler.Handle)

	r.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		log.Printf("Panic occurred: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": fmt.Sprintf("%v", err),
		})
	}))

	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")
}
