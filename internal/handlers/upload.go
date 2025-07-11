package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadHandler обрабатывает загрузку файлов
type UploadHandler struct {
	uploadsDir string
}

// NewUploadHandler создает новый экземпляр UploadHandler
func NewUploadHandler(uploadsDir string) *UploadHandler {
	return &UploadHandler{uploadsDir: uploadsDir}
}

// Handle обрабатывает POST запросы для загрузки файлов
func (h *UploadHandler) Handle(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерация уникального имени файла
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(h.uploadsDir, newFilename)

	// Сохранение файла
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    newFilename,
		"image": "/static/uploads/" + newFilename,
	})
}
