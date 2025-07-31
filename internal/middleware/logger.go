package middleware

import (
	"subscription-aggregator/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()

		logger.Log.Infof("Started %s %s for %s", context.Request.Method, context.Request.URL.Path, context.ClientIP())

		context.Next()

		duration := time.Since(start)
		statusCode := context.Writer.Status()
		method := context.Request.Method
		path := context.Request.URL.Path
		clientIP := context.ClientIP()

		if len(context.Errors) > 0 {
			for _, e := range context.Errors.Errors() {
				logger.Log.Errorf("Request error: %s", e)
			}
		}

		logger.Log.WithFields(map[string]interface{}{
			"status":   statusCode,
			"method":   method,
			"path":     path,
			"duration": duration,
			"clientIP": clientIP,
		}).Info("completed request")
	}
}
