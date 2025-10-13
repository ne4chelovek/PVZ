package handler

import (
	"PVZ/internal/metrics"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	product, err := h.receptionService.AddProduct(ctx, req.PVZID, req.Type)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timeout"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	metrics.ProductsAddedTotal.Inc()

	c.JSON(http.StatusCreated, product)
}
