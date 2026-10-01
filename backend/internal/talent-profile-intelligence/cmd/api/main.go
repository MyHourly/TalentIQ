package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"talentiq/talent-profile-intelligence/internal/config"
	"talentiq/talent-profile-intelligence/internal/database"
	"talentiq/talent-profile-intelligence/internal/handler"
	"talentiq/talent-profile-intelligence/internal/logger"
	"talentiq/talent-profile-intelligence/internal/messaging"
	"talentiq/talent-profile-intelligence/internal/repository"
	"talentiq/talent-profile-intelligence/internal/service"
)

func main() {

	// ---------------------------------------------------------
	// Configuration
	// ---------------------------------------------------------

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

	// ---------------------------------------------------------
	// Database
	// ---------------------------------------------------------

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

	// ---------------------------------------------------------
	// Kafka Producer
	// ---------------------------------------------------------

	// Create Kafka event publisher.
	//
	// DEVELOPMENT ONLY:
	// Kafka currently runs locally through Docker.
	//
	// FUTURE REFACTOR:
	// Broker and topic configuration should eventually come
	// from the centralized TalentIQ Kafka infrastructure.
	//
	// The service layer depends only on the EventPublisher
	// interface, so the Kafka implementation can be replaced
	// later without changing business logic.
	eventPublisher := messaging.NewKafkaProducer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		appLogger.Logger,
	)

	// Close Kafka producer when the application exits.
	defer func() {

		if err := eventPublisher.Close(); err != nil {

			appLogger.Error(
				"failed to close kafka producer",
				"error", err,
			)
		}
	}()

	appLogger.Info(
		"Kafka event publisher initialized",
		"brokers", cfg.KafkaBrokers,
		"topic", cfg.KafkaTopic,
	)

	// ---------------------------------------------------------
	// Kafka Consumer
	// ---------------------------------------------------------

	// Create Kafka consumer.
	//
	// DEVELOPMENT ONLY:
	// Kafka currently runs locally through Docker.
	//
	// FUTURE REFACTOR:
	// Kafka broker, topic and consumer-group configuration
	// should eventually come from the centralized TalentIQ
	// Kafka infrastructure.
	kafkaConsumer := messaging.NewKafkaConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		cfg.KafkaGroupID,
		appLogger.Logger,
	)

	// Close Kafka consumer when the application exits.
	defer func() {

		if err := kafkaConsumer.Close(); err != nil {

			appLogger.Error(
				"failed to close kafka consumer",
				"error", err,
			)
		}
	}()

	// Start Kafka consumer in the background.
	//
	// The Kafka consumer runs independently from the HTTP API.
	go func() {

		if err := kafkaConsumer.Consume(ctx); err != nil {

			appLogger.Error(
				"Kafka consumer stopped",
				"error", err,
			)
		}

	}()

	// ---------------------------------------------------------
	// Repository
	// ---------------------------------------------------------

	// Create Talent Profile repository.
	talentProfileRepository :=
		repository.NewTalentProfileRepository(
			db,
			appLogger.Logger,
		)

	// ---------------------------------------------------------
	// Service
	// ---------------------------------------------------------

	// Create Talent Profile service.
	//
	// The service receives EventPublisher instead of directly
	// depending on KafkaProducer.
	//
	// This keeps business logic independent from Kafka.
	//
	// FUTURE REFACTOR:
	// A centralized TalentIQ messaging implementation can be
	// injected here without changing the service layer.
	talentProfileService :=
		service.NewTalentProfileService(
			talentProfileRepository,
			eventPublisher,
			appLogger.Logger,
		)

	// ---------------------------------------------------------
	// HTTP Handler
	// ---------------------------------------------------------

	// Create Talent Profile HTTP handler.
	talentProfileHandler :=
		handler.NewTalentProfileHandler(
			talentProfileService,
			appLogger.Logger,
		)

	// ---------------------------------------------------------
	// Gin Router
	// ---------------------------------------------------------

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

	// ---------------------------------------------------------
	// Versioned Talent Profile API
	// ---------------------------------------------------------

	v1 := router.Group("/api/v1")
	{
		talentProfiles := v1.Group("/talent-profiles")
		{
			// Create talent profile.
			talentProfiles.POST(
				"",
				talentProfileHandler.Create,
			)

			// List talent profiles.
			talentProfiles.GET(
				"",
				talentProfileHandler.List,
			)

			// Get talent profile by ID.
			talentProfiles.GET(
				"/:id",
				talentProfileHandler.GetByID,
			)

			// Update talent profile.
			talentProfiles.PUT(
				"/:id",
				talentProfileHandler.Update,
			)

			// Deactivate talent profile.
			talentProfiles.DELETE(
				"/:id",
				talentProfileHandler.Delete,
			)
		}
	}

	// ---------------------------------------------------------
	// Start HTTP Server
	// ---------------------------------------------------------

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
