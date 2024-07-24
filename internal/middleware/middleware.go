package middleware

import (
	"context"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

func AccessLog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		ctx := context.WithValue(c.Request.Context(), "requestID", requestID)
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()

		c.Next()

		logger.Info("New request",
			"requestID", requestID,
			"method", c.Request.Method,
			"remote_addr", c.Request.RemoteAddr,
			"url", c.Request.URL.Path,
			"time", time.Since(start),
		)
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(-1, c.Errors.JSON())
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func GzipMiddleware() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}

func RateLimitingMiddleware() gin.HandlerFunc {
	limiter := tollbooth.NewLimiter(10, nil) // 1 request per second
	return tollbooth_gin.LimitHandler(limiter)
}

/* ValidateRequest middleware validates the request body
func ValidateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var json MyStruct
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}*/
