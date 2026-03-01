package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"Authentication_Service/internal/dto/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GlobalPanicHandler() gin.RecoveryFunc {
	return func(c *gin.Context, err any) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse[any]{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}
}

// RequestIDMiddleware adds a unique request ID to each request for tracking
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate a simple unique ID using timestamp and random bytes
			b := make([]byte, 4)
			rand.Read(b)
			requestID = fmt.Sprintf("%d-%s", time.Now().UnixNano(), hex.EncodeToString(b))
		}
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// RequestLogger logs request details with request ID for debugging duplicate requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID, exists := c.Get("RequestID")
		if !exists {
			requestID = "unknown"
		}
		
		c.Next()
		
		latency := time.Since(start)
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		
		fmt.Fprintf(gin.DefaultWriter, "[GIN] %s | %3d | %13v | %15s | %-7s %s | RequestID: %v\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			c.Request.Method,
			c.Request.RequestURI,
			requestID,
		)
	}
}
