package middleware

import (
	"subscription-aggregator/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()

		context.Next()

		duration := time.Since(start)
		statusCode := context.Writer.Status()
		method := context.Request.Method
		path := context.Request.URL.Path
		clientIP := context.ClientIP()

		logger.Log.WithFields(map[string]interface{}{
			"status":   statusCode,
			"method":   method,
			"path":     path,
			"duration": duration,
			"clientIP": clientIP,
		}).Info("incoming request")
	}
}
