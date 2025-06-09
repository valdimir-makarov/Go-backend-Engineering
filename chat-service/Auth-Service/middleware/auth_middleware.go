package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Security headers middleware
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Writer.Header().Set("Referrer-Policy", "no-referrer")
		c.Next()
	}
}

// Request tracing middleware
func TraceRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("requestID", requestID)

		start := time.Now()
		log.Printf("Started %s %s | RequestID: %s", c.Request.Method, c.Request.URL.Path, requestID)

		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		log.Printf("Completed %s %s | Status: %d | Duration: %v | RequestID: %s",
			c.Request.Method, c.Request.URL.Path, status, duration, requestID)
	}
}
