package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"talentiq/talent-profile-intelligence/internal/config"
	"talentiq/talent-profile-intelligence/internal/database"
	"talentiq/talent-profile-intelligence/internal/handler"
	"talentiq/talent-profile-intelligence/internal/logger"
	"talentiq/talent-profile-intelligence/internal/repository"
	"talentiq/talent-profile-intelligence/internal/service"
)

func main() {
	// Load application configuration.
	cfg := config.Load()

	// Create application logger.
	appLogger := logger.New(cfg.AppEnv)

	appLogger.Info(
		"starting talent profile intelligence service",
		"environment", cfg.AppEnv,
		"port", cfg.ServerPort,
	)

	// Root context used during application startup.
	ctx := context.Background()

	// Create PostgreSQL connection pool.
	db, err := database.NewPostgres(ctx, cfg)
	if err != nil {
		appLogger.Error(
			"database initialization failed",
			"error", err,
		)

		log.Fatal(err)
	}

	// Close database connections when application exits.
	defer db.Close()

	appLogger.Info(
		"PostgreSQL connection established",
	)

	// Create repository.
	talentProfileRepository :=
		repository.NewTalentProfileRepository(
			db,
			appLogger.Logger,
		)

	// Create service.
	talentProfileService :=
		service.NewTalentProfileService(
			talentProfileRepository,
			appLogger.Logger,
		)

	// The service will be injected into HTTP handlers
	// in the next lesson.
	_ = talentProfileService

	// Create Gin router.
	router := gin.Default()

	// Health endpoint.
	router.GET(
		"/health",
		handler.Health,
	)

	// Readiness endpoint.
	router.GET(
		"/ready",
		handler.Ready(db),
	)

	appLogger.Info(
		"HTTP server started",
		"port", cfg.ServerPort,
	)

	// Start HTTP server.
	if err := router.Run(
		":" + cfg.ServerPort,
	); err != nil {

		appLogger.Error(
			"HTTP server stopped with error",
			"error", err,
		)

		log.Fatal(err)
	}
}
