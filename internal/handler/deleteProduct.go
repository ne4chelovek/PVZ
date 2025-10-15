package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *PVZHandler) DeleteLastProduct(c *gin.Context) {
	pvzID := c.Param("pvzId")
	if pvzID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pvzId is required"})
		return
	}

	err := h.receptionService.DeleteLastProduct(c.Request.Context(), pvzID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
