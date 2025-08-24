package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *PVZHandler) CloseLastReception(c *gin.Context) {
	pvzID := c.Param("pvzId")
	if pvzID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pvzId is required"})
		return
	}

	reception, err := h.receptionService.CloseLastReception(c.Request.Context(), pvzID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reception)
}
