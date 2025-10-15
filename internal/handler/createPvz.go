package handler

import (
	"PVZ/internal/logger"
	"PVZ/internal/metrics"
	"PVZ/internal/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (h *PVZHandler) CreatePVZ(c *gin.Context) {
	var input struct {
		ID               string `json:"id"`
		RegistrationDate string `json:"registrationDate"`
		City             string `json:"city"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	pvz := model.PVZ{
		ID:   input.ID,
		City: input.City,
	}

	createdPVZ, err := h.pvzService.CreatePVZ(c.Request.Context(), &pvz)
	if err != nil {
		logger.Error("Failed to create PVZ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PVZ"})
		return
	}

	metrics.PVZCreatedTotal.Inc()

	c.JSON(http.StatusCreated, createdPVZ)
}
