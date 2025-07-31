package utils

import (
	"net/http"
	"subscription-aggregator/pkg/logger"

	"github.com/gin-gonic/gin"
)

func BindJSONOrAbort[T any](context *gin.Context, target *T) bool {
	logger.Log.Infof("Attempting to bind JSON to %T", target)
	if err := context.ShouldBindJSON(target); err != nil {
		logger.Log.Errorf("Failed to bind JSON: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON: " + err.Error(),
		})
		return false
	}
	logger.Log.Infof("Successfully bound JSON to %T", target)
	return true
}
