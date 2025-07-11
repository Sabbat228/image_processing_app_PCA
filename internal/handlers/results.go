package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ResultsHandler обрабатывает запросы к результатам
type ResultsHandler struct {
	uploadsDir string
}

// NewResultsHandler создает новый экземпляр ResultsHandler
func NewResultsHandler(uploadsDir string) *ResultsHandler {
	return &ResultsHandler{uploadsDir: uploadsDir}
}

// Handle обрабатывает GET запросы для просмотра результатов
func (h *ResultsHandler) Handle(c *gin.Context) {
	id := c.Param("id")
	filePath := filepath.Join(h.uploadsDir, id)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Image not found",
		})
		return
	}

	c.HTML(http.StatusOK, "results.html", gin.H{
		"id":    id,
		"image": "/static/uploads/" + id,
	})
}
