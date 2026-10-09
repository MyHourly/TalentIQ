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
	"talentiq/talent-profile-intelligence/internal/messaging/handlers"
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
	eventPublisher := messaging.NewKafkaProducer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		appLogger.Logger,
	)

	// Close Kafka producer when application exits.
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
	// Repository
	// ---------------------------------------------------------

	// Create Talent Profile repository.
	talentProfileRepository :=
		repository.NewTalentProfileRepository(
			db,
			appLogger.Logger,
		)

	// Create Talent Preference repository.
	talentPreferenceRepository :=
		repository.NewTalentPreferenceRepository(
			db,
			appLogger.Logger,
		)

	// TPI-015: Create Talent Skill repository.
	talentSkillRepository :=
		repository.NewTalentSkillRepository(
			db,
			appLogger.Logger,
		)

	// TPI-016: Create Talent Certification repository.
	talentCertificationRepository :=
		repository.NewTalentCertificationRepository(
			db,
			appLogger.Logger,
		)

	// TPI-017: Create Talent Project repository.
	talentProjectRepository :=
		repository.NewTalentProjectRepository(
			db,
			appLogger.Logger,
		)

	// TPI-018: Create Talent Experience repository.
	talentExperienceRepository :=
		repository.NewTalentExperienceRepository(
			db,
			appLogger.Logger,
		)

	// ---------------------------------------------------------
	// Event Processing Service
	// ---------------------------------------------------------

	// Create the service responsible for business processing
	// triggered by Kafka talent profile events.
	talentProfileEventService :=
		service.NewTalentProfileEventService(
			talentProfileRepository,
			appLogger.Logger,
		)

	// ---------------------------------------------------------
	// Kafka Event Dispatcher
	// ---------------------------------------------------------

	// Create event dispatcher.
	eventDispatcher := messaging.NewEventDispatcher(
		appLogger.Logger,
	)

	// ---------------------------------------------------------
	// Profile Created Handler
	// ---------------------------------------------------------

	// The created handler uses the event processing service
	// to perform business processing.
	profileCreatedHandler :=
		handlers.NewProfileCreatedHandler(
			talentProfileEventService,
			appLogger.Logger,
		)

	eventDispatcher.RegisterHandler(
		"candidate.profile.created",
		profileCreatedHandler,
	)

	// ---------------------------------------------------------
	// Profile Updated Handler
	// ---------------------------------------------------------

	// The updated handler also uses the event processing service.
	profileUpdatedHandler :=
		handlers.NewProfileUpdatedHandler(
			talentProfileEventService,
			appLogger.Logger,
		)

	eventDispatcher.RegisterHandler(
		"candidate.profile.updated",
		profileUpdatedHandler,
	)

	// ---------------------------------------------------------
	// Profile Deactivated Handler
	// ---------------------------------------------------------

	// The deactivated handler uses the event processing service
	// to process the deactivated profile event.
	profileDeactivatedHandler :=
		handlers.NewProfileDeactivatedHandler(
			talentProfileEventService,
			appLogger.Logger,
		)

	eventDispatcher.RegisterHandler(
		"candidate.profile.deactivated",
		profileDeactivatedHandler,
	)

	// ---------------------------------------------------------
	// Kafka Consumer
	// ---------------------------------------------------------

	// Create Kafka consumer.
	//
	// The consumer receives the EventDispatcher so that
	// Kafka messages can be routed to the correct handler.
	kafkaConsumer := messaging.NewKafkaConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		cfg.KafkaGroupID,
		appLogger.Logger,
		eventDispatcher,
	)

	// Close Kafka consumer when application exits.
	defer func() {
		if err := kafkaConsumer.Close(); err != nil {
			appLogger.Error(
				"failed to close kafka consumer",
				"error", err,
			)
		}
	}()

	// ---------------------------------------------------------
	// Start Kafka Consumer
	// ---------------------------------------------------------

	// Start Kafka consumer in the background.
	//
	// Event flow:
	//
	// Kafka Topic
	//      ↓
	// Kafka Consumer
	//      ↓
	// Event Dispatcher
	//      ↓
	// Event Handler
	//      ↓
	// Event Processing Service
	go func() {
		if err := kafkaConsumer.Consume(ctx); err != nil {
			appLogger.Error(
				"Kafka consumer stopped",
				"error", err,
			)
		}
	}()

	// ---------------------------------------------------------
	// Talent Profile Service
	// ---------------------------------------------------------

	// Create Talent Profile service.
	//
	// The service depends on EventPublisher instead of
	// directly depending on KafkaProducer.
	talentProfileService :=
		service.NewTalentProfileService(
			talentProfileRepository,
			eventPublisher,
			appLogger.Logger,
		)

	// ---------------------------------------------------------
	// Talent Preference Service
	// ---------------------------------------------------------

	// Create Talent Preference service.
	//
	// This service handles the business logic for
	// talent preferences and availability.
	talentPreferenceService :=
		service.NewTalentPreferenceService(
			talentPreferenceRepository,
			appLogger.Logger,
		)

	// ---------------------------------------------------------
	// Talent Skill Service
	// ---------------------------------------------------------

	// TPI-015: Create Talent Skill service.
	//
	// This service handles the business rules for adding,
	// updating and removing a talent's skills.
	talentSkillService :=
		service.NewTalentSkillService(
			talentSkillRepository,
			appLogger.Logger,
		)

	// ---------------------------------------------------------
	// Talent Certification Service
	// ---------------------------------------------------------

	// TPI-016: Create Talent Certification service.
	//
	// This service validates certification details and
	// decides the certification status (ACTIVE / EXPIRED).
	talentCertificationService :=
		service.NewTalentCertificationService(
			talentCertificationRepository,
			appLogger.Logger,
		)

	// ---------------------------------------------------------
	// Talent Project Service
	// ---------------------------------------------------------

	// TPI-017: Create Talent Project service.
	//
	// This service validates project details, adds a project
	// for a talent, and marks it completed.
	talentProjectService :=
		service.NewTalentProjectService(
			talentProjectRepository,
			appLogger.Logger,
		)

	// ---------------------------------------------------------
	// Talent Experience Service
	// ---------------------------------------------------------

	// TPI-018: Create Talent Experience service.
	//
	// This service validates a job (company, title, dates)
	// before it is saved in the talent's work history.
	talentExperienceService :=
		service.NewTalentExperienceService(
			talentExperienceRepository,
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

	// Create Talent Preference HTTP handler.
	//
	// The same handler also handles the Talent Availability
	// endpoint because availability is stored in the existing
	// talent_preferences table.
	talentPreferenceHandler :=
		handler.NewTalentPreferenceHandler(
			talentPreferenceService,
			appLogger.Logger,
		)

	// TPI-015: Create Talent Skill HTTP handler.
	talentSkillHandler :=
		handler.NewTalentSkillHandler(
			talentSkillService,
			appLogger.Logger,
		)

	// TPI-016: Create Talent Certification HTTP handler.
	talentCertificationHandler :=
		handler.NewTalentCertificationHandler(
			talentCertificationService,
			appLogger.Logger,
		)

	// TPI-017: Create Talent Project HTTP handler.
	talentProjectHandler :=
		handler.NewTalentProjectHandler(
			talentProjectService,
			appLogger.Logger,
		)

	// TPI-018: Create Talent Experience HTTP handler.
	talentExperienceHandler :=
		handler.NewTalentExperienceHandler(
			talentExperienceService,
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
	// Versioned APIs
	// ---------------------------------------------------------

	v1 := router.Group("/api/v1")
	{

		// -----------------------------------------------------
		// Talent Profile APIs
		// -----------------------------------------------------

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

			// -------------------------------------------------
			// Talent Preference APIs
			// -------------------------------------------------

			// Create talent preferences.
			talentProfiles.POST(
				"/:id/preferences",
				talentPreferenceHandler.Create,
			)

			// Get talent preferences.
			talentProfiles.GET(
				"/:id/preferences",
				talentPreferenceHandler.GetByTalentID,
			)

			// Update talent preferences.
			talentProfiles.PUT(
				"/:id/preferences",
				talentPreferenceHandler.Update,
			)
		}

		// -----------------------------------------------------
		// Talent APIs (client work-plan paths: /api/v1/talents/...)
		// -----------------------------------------------------

		talents := v1.Group("/talents")
		{

			// Talent Availability API.
			//
			// The availability endpoint follows the client work-plan
			// path: PATCH /api/v1/talents/:id/availability
			//
			// Availability is stored in the existing
			// talent_preferences.availability_status field.
			talents.PATCH(
				"/:id/availability",
				talentPreferenceHandler.UpdateAvailability,
			)

			// TPI-015: Add a skill to a talent, or change its
			// proficiency level if the talent already has that skill.
			talents.PUT(
				"/:id/skills/:skillId",
				talentSkillHandler.Upsert,
			)

			// TPI-015: Remove a skill from a talent.
			talents.DELETE(
				"/:id/skills/:skillId",
				talentSkillHandler.Remove,
			)

			// TPI-016: Add a certification to a talent.
			talents.POST(
				"/:id/certifications",
				talentCertificationHandler.Add,
			)

			// TPI-017: Add a project to a talent.
			talents.POST(
				"/:id/projects",
				talentProjectHandler.Add,
			)

			// TPI-017: Mark one of the talent's projects as completed.
			//
			// NOTE: the wildcard names (:id, :projectId) must match
			// across routes at the same position, or Gin panics at startup.
			talents.POST(
				"/:id/projects/:projectId/complete",
				talentProjectHandler.Complete,
			)

			// TPI-018: Add a job to a talent's work history.
			talents.POST(
				"/:id/experience",
				talentExperienceHandler.Add,
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
