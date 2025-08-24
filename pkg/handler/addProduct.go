package handler

import (
	"PVZ/internal/metrics"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *PVZHandler) AddProduct(c *gin.Context) {
	var req struct {
		Type  string `json:"type" binding:"required,oneof=электроника одежда обувь"`
		PVZID string `json:"pvzId" binding:"required,uuid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
		return
	}

	product, err := h.receptionService.AddProduct(c.Request.Context(), req.PVZID, req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metrics.ProductsAddedTotal.Inc()

	c.JSON(http.StatusCreated, product)
}
