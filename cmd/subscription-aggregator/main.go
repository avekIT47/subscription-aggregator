package main

import (
	"subscription-aggregator/internal/config"
	"subscription-aggregator/internal/handler"
	"subscription-aggregator/internal/middleware"
	"subscription-aggregator/internal/repository"
	"subscription-aggregator/internal/service"
	"subscription-aggregator/migrations"
	"subscription-aggregator/pkg/database"
	"subscription-aggregator/pkg/logger"

	_ "subscription-aggregator/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Subscription Aggregator API
// @version 1.0
// @description Это API для управления подписками
// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	cfg := config.LoadConfig("config/config.yaml")
	dsn := cfg.GetDSN()

	db := database.InitDB(dsn)

	logger.InitLogger()

	if err := migrations.AutoMigrate(db); err != nil {
		logger.Log.Fatalf("transfer start error: %v", err)
	}

	subRepo := repository.NewSubscriptionRepository(db)
	subService := service.NewSubscriptionService(subRepo)
	subHandler := handler.NewSubscriptionHandler(subService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{
		sub := api.Group("/subscriptions")
		{
			sub.POST("/", subHandler.CreateSubscriprion)
			sub.GET("/:id", subHandler.GetSubscriptionByID)
			sub.PUT("/:id", subHandler.UpdateSubscription)
			sub.DELETE("/:id", subHandler.DeleteSubscription)
			sub.GET("/user/:id/total", subHandler.GetTotal)
			sub.POST("/:id/list", subHandler.GetSubscriptionsList)
		}
	}

	logger.Log.Info("starting the server. Port :8080")
	if err := router.Run(":8080"); err != nil {
		logger.Log.Fatalf("server startup error: %v", err)
	}
}
