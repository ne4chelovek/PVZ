package handler

import (
	"PVZ/internal/metrics"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *PVZHandler) CreateReception(c *gin.Context) {
	var req struct {
		PVZID string `json:"pvzId" binding:"required,uuid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	reception, err := h.receptionService.CreateReception(c.Request.Context(), req.PVZID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metrics.ReceptionsCreatedTotal.Inc()

	c.JSON(http.StatusCreated, reception)
}
