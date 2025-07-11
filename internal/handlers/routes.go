package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image-processing-app/internal/services"
	"log"
	"net/http"
)

func SetupRoutes(r *gin.Engine) {
	// Инициализация сервисов
	imageProcessor := services.NewImageProcessor(50, 200)
	uploadsDir := "web/static/uploads"

	// Инициализация обработчиков (теперь с экспортируемыми типами)
	homeHandler := NewHomeHandler()
	uploadHandler := NewUploadHandler(uploadsDir)
	resultsHandler := NewResultsHandler(uploadsDir)
	analyzeHandler := NewAnalyzeHandler(imageProcessor, uploadsDir)

	// Настройка маршрутов
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

	// Статические файлы
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")
}
