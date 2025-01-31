package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/yakob-abada/go-rest-user/config"
	_ "github.com/yakob-abada/go-rest-user/docs"
	"github.com/yakob-abada/go-rest-user/pkg/bootstrap"
	"github.com/yakob-abada/go-rest-user/pkg/logging"
)

func main() {
	e := echo.New()
	db := config.InitDB()

	rabbitMQ, _ := config.InitRabbitMQ()
	defer rabbitMQ.Close()

	logger := logging.NewZeroLogger()
	userHandler := bootstrap.NewUserHandler(db, rabbitMQ, logger)

	e.PUT("/users/:id", userHandler.UpdateUser)
	e.POST("/users", userHandler.SaveUser)
	e.GET("/users", userHandler.GetUsers)
	e.DELETE("/users/:id", userHandler.DeleteUser)

	e.GET("/health", bootstrap.NewHealthHandler(db).HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	ctx := context.Background()
	// Start server
	go func() {
		logger.Info(ctx, "Starting server", map[string]interface{}{"port": "8080"})
		if err := e.Start(":8080"); err != nil {
			logger.Error(ctx, "Server failed", map[string]interface{}{"error": err.Error()})
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Warn(ctx, "Shutting down server", nil)

	if err := e.Shutdown(context.Background()); err != nil {
		logger.Error(ctx, "Server forced to shutdown", map[string]interface{}{"error": err.Error()})
	}

	logger.Info(ctx, "Server stopped cleanly", nil)
}
