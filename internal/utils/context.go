package utils

import (
	"net/http"
	"subscription-aggregator/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckID(context *gin.Context) (uuid.UUID, bool) {
	idParam := context.Param("id")
	logger.Log.Infof("Parsing UUID from param: %s", idParam)
	id, err := uuid.Parse(idParam)
	if err != nil {
		logger.Log.Errorf("Invalid UUID format: %s, error: %v", idParam, err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return uuid.Nil, false
	}
	logger.Log.Infof("Successfully parsed UUID: %s", id)
	return id, true
}

func GetDate(context *gin.Context) (time.Time, *time.Time) {
	fromStr := context.Query("from")
	toStr := context.Query("to")
	logger.Log.Infof("Parsing dates from query params: from='%s', to='%s'", fromStr, toStr)

	var from time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			logger.Log.Errorf("Invalid 'from' date format: %s, error: %v", fromStr, err)
			context.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' date"})
			return time.Time{}, nil
		}
		logger.Log.Infof("Parsed 'from' date: %s", from.Format("2006-01-02"))
	}
	var to *time.Time
	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			logger.Log.Errorf("Invalid 'to' date format: %s, error: %v", toStr, err)
			context.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' date"})
			return time.Time{}, nil
		}
		to = &t
		logger.Log.Infof("Parsed 'to' date: %s", to.Format("2006-01-02"))
	}

	return from, to
}
