package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler обрабатывает запросы к главной странице
type HomeHandler struct{}

// NewHomeHandler создает новый экземпляр HomeHandler
func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// Handle обрабатывает GET запрос к корневому маршруту
func (h *HomeHandler) Handle(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
