package deploy

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"subscription-aggregator/pkg/logger"
	"time"
)

func RunServer(handler http.Handler, port string) error {
	srv := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	go func() {
		logger.Log.Infof("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Listen error: %v", err)
		}
		logger.Log.Info("ListenAndServe exited")
	}()

	logger.Log.Info("Server is running. Waiting for interrupt signal to shut down...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit

	logger.Log.Infof("Received signal: %v", sig)
	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Errorf("Server shutdown error: %v", err)
		return err
	}

	logger.Log.Info("Server exited gracefully")
	return nil
}
