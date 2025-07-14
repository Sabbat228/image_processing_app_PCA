package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"image-processing-app/internal/services"
	"image-processing-app/internal/utils"
)

type AnalyzeHandler struct {
	imageProcessor *services.ImageProcessor
	uploadsDir     string
}

func NewAnalyzeHandler(ip *services.ImageProcessor, uploadsDir string) *AnalyzeHandler {
	return &AnalyzeHandler{
		imageProcessor: ip,
		uploadsDir:     uploadsDir,
	}
}

func (h *AnalyzeHandler) Handle(c *gin.Context) {
	// Define request structure
	var request struct {
		ImageID  string `json:"image_id"`
		Method   string `json:"method" binding:"required,oneof=pca nmf"`
		NFactors int    `json:"n_factors" binding:"required,min=1,max=100"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request parameters",
			"details": err.Error(),
		})
		return
	}

	if request.NFactors < 1 || request.NFactors > 100 {
		request.NFactors = 10 // Значение по умолчанию
	}

	imgPath := filepath.Join(h.uploadsDir, request.ImageID)
	img, err := utils.LoadImage(imgPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to load image",
			"details": err.Error(),
		})
		return
	}

	if img.Bounds().Empty() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Empty image provided",
		})
		return
	}

	result, err := h.imageProcessor.ProcessImage(request.Method, img, request.NFactors)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Image processing failed",
			"details": err.Error(),
		})
		return
	}

	resultFilename := fmt.Sprintf("processed_%s_%d_%s",
		request.Method,
		request.NFactors,
		request.ImageID)
	resultPath := filepath.Join(h.uploadsDir, resultFilename)

	if err := utils.SaveImage(resultPath, result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save result",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":  "/static/uploads/" + resultFilename,
		"method":  request.Method,
		"factors": request.NFactors,
	})
}
