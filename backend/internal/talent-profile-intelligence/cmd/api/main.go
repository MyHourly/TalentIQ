package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"talentiq/talent-profile-intelligence/internal/config"
	"talentiq/talent-profile-intelligence/internal/database"
	"talentiq/talent-profile-intelligence/internal/handler"
	"talentiq/talent-profile-intelligence/internal/logger"
)

func main() {
	// Load application configuration.
	cfg := config.Load()

	// Create the application logger.
	appLogger := logger.New(cfg.AppEnv)

	appLogger.Info(
		"starting talent profile intelligence service",
		"environment", cfg.AppEnv,
		"port", cfg.ServerPort,
	)

	// Create a root context for application startup.
	ctx := context.Background()

	// Connect to PostgreSQL.
	db, err := database.NewPostgres(ctx, cfg)
	if err != nil {
		appLogger.Error(
			"database initialization failed",
			"error", err,
		)

		// Exit because the service cannot operate without
		// its required database.
		log.Fatal(err)
	}

	// Close the database connection pool when the
	// application stops.
	defer db.Close()

	appLogger.Info("PostgreSQL connection established")

	// Create the Gin router.
	router := gin.Default()

	// Basic application health check.
	router.GET("/health", handler.Health)

	// Database readiness check.
	router.GET("/ready", handler.Ready(db))

	appLogger.Info(
		"HTTP server started",
		"port", cfg.ServerPort,
	)

	// Start the HTTP server.
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		appLogger.Error(
			"HTTP server stopped with error",
			"error", err,
		)

		log.Fatal(err)
	}
}