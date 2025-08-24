package middleware

import (
	"fmt"
	"time"

	"PVZ/internal/metrics"
	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := c.Writer.Status()
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path // fallback
		}

		metrics.HttpRequestTotal.WithLabelValues(
			c.Request.Method,
			endpoint,
			fmt.Sprintf("%d", status),
		).Inc()

		metrics.HttpRequestDuration.WithLabelValues(
			c.Request.Method,
			endpoint,
		).Observe(time.Since(start).Seconds())
	}
}
