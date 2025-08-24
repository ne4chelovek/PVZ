package handler

import (
	"PVZ/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (h *PVZHandler) GetAllPVZ(c *gin.Context) {
	var startDate, endDate *string
	if s := c.Query("startDate"); s != "" {
		startDate = &s
	}
	if e := c.Query("endDate"); e != "" {
		endDate = &e
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	pvzs, err := h.pvzService.GetAllPVZWithReceptions(c.Request.Context(), startDate, endDate, page, limit)
	if err != nil {
		logger.Error("Failed to get PVZ list", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PVZ list"})
		return
	}

	c.JSON(http.StatusOK, pvzs)
}
