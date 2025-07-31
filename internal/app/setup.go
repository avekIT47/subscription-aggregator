package app

import (
	"strings"
	"subscription-aggregator/internal/config"
	"subscription-aggregator/internal/handler"
	"subscription-aggregator/internal/middleware"
	"subscription-aggregator/internal/repository"
	"subscription-aggregator/internal/service"
	"subscription-aggregator/migrations"
	"subscription-aggregator/pkg/database"
	"subscription-aggregator/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup() (router *gin.Engine, port string, dbCloser func() error, err error) {
	cfg := config.LoadConfig("config/config.yaml")
	logger.Log.Info("Config loaded")

	dsn := cfg.GetDSN()

	port = cfg.Server.Port
	if port == "" {
		port = ":8080"
	} else if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	logger.Log.Infof("Server will start on port %s", port)

	db := database.InitDB(dsn)
	logger.Log.Info("Database connection initialized")

	if err := migrations.AutoMigrate(db); err != nil {
		logger.Log.Fatalf("migration error: %v", err)
	}

	subRepo := repository.NewSubscriptionRepository(db)
	subService := service.NewSubscriptionService(subRepo)
	subHandler := handler.NewSubscriptionHandler(subService)

	router = gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		sub := api.Group("/subscriptions")
		{
			sub.POST("/:user_id", subHandler.CreateSubscriprion)
			sub.GET("/:id", subHandler.GetSubscriptionByID)
			sub.PUT("/:id", subHandler.UpdateSubscription)
			sub.DELETE("/:id", subHandler.DeleteSubscription)
			sub.GET("/user/:user_id/total", subHandler.GetTotal)
			sub.POST("/:user_id/list", subHandler.GetSubscriptionsList)
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, "", nil, err
	}

	dbCloser = func() error {
		logger.Log.Info("Closing database connection")
		return sqlDB.Close()
	}

	logger.Log.Info("Setup finished successfully")

	return router, port, dbCloser, nil
}
