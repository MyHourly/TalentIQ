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
	// Create service.
// Create service.
talentProfileService :=
    service.NewTalentProfileService(
        talentProfileRepository,
        nil,
        appLogger.Logger,
    )

	// Create Talent Profile HTTP handler.
	talentProfileHandler :=
		handler.NewTalentProfileHandler(
			talentProfileService,
			appLogger.Logger,
		)

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

	// Versioned Talent Profile API.
	v1 := router.Group("/api/v1")
	{
		talentProfiles := v1.Group("/talent-profiles")
		{
			talentProfiles.POST(
				"",
				talentProfileHandler.Create,
			)

			talentProfiles.GET(
				"",
				talentProfileHandler.List,
			)

			talentProfiles.GET(
				"/:id",
				talentProfileHandler.GetByID,
			)

			talentProfiles.PUT(
				"/:id",
				talentProfileHandler.Update,
			)

			talentProfiles.DELETE(
				"/:id",
				talentProfileHandler.Delete,
			)
		}
	}

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
