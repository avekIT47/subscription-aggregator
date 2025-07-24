package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckID(context *gin.Context) (uuid.UUID, bool) {
	idParam := context.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return uuid.Nil, false
	}
	return id, true
}
