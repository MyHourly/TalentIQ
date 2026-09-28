package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"talentiq/talent-profile-intelligence/internal/config"
	"talentiq/talent-profile-intelligence/internal/database"
	"talentiq/talent-profile-intelligence/internal/handler"
)

func main() {
	// Load application configuration from environment variables.
	cfg := config.Load()

	// Create a root context for application startup.
	ctx := context.Background()

	// Connect to PostgreSQL before starting the HTTP server.
	db, err := database.NewPostgres(ctx, cfg)
	if err != nil {
		log.Fatalf(
			"database initialization failed: %v",
			err,
		)
	}

	// Close the database connection pool when the application exits.
	defer db.Close()

	log.Println("PostgreSQL connection established")

	// Create the Gin HTTP router.
	router := gin.Default()

	// Liveness endpoint.
	//
	// This checks whether the application itself is running.
	router.GET("/health", handler.Health)

	// Readiness endpoint.
	//
	// This checks whether the application can communicate
	// with PostgreSQL.
	router.GET("/ready", handler.Ready(db))

	// Start the HTTP server.
	log.Printf(
		"%s starting on port %s",
		cfg.AppName,
		cfg.ServerPort,
	)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf(
			"failed to start HTTP server: %v",
			err,
		)
	}
}