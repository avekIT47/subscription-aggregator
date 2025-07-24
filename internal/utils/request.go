package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindJSONOrAbort[T any](context *gin.Context, target *T) bool {
	if err := context.ShouldBindJSON(target); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON: " + err.Error(),
		})
		return false
	}
	return true
}
