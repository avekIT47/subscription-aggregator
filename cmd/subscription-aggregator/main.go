package main

import (
	"subscription-aggregator/internal/app"
	"subscription-aggregator/internal/deploy"
	"subscription-aggregator/pkg/logger"

	_ "subscription-aggregator/docs"
)

// @title Subscription Aggregator API
// @version 1.0
// @description Это API для управления подписками
// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	logger.InitLogger()
	logger.Log.Info("Starting application setup...")

	router, port, dbCloser, err := app.Setup()
	if err != nil {
		logger.Log.Fatalf("Setup error: %v", err)
	}
	logger.Log.Infof("Setup completed successfully. Server will start on port %s", port)

	if err := deploy.RunServer(router, port); err != nil {
		logger.Log.Fatalf("Server error: %v", err)
	}
	logger.Log.Info("Server has stopped running.")

	logger.Log.Info("Closing database connection...")
	if err := dbCloser(); err != nil {
		logger.Log.Errorf("DB close error: %v", err)
	} else {
		logger.Log.Info("Database connection closed successfully.")
	}

	logger.Log.Info("Application shutdown complete.")
}
