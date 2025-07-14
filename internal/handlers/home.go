package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Handle(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
