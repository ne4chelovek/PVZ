package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *PVZHandler) LoginDummy(c *gin.Context) {
	var req struct {
		Role string `json:"role" binding:"required,oneof=employee moderator"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing role"})
		return
	}

	token, err := h.authService.LoginDummy(c.Request.Context(), req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
