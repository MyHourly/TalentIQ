package config

import "os"

// Config contains all configuration required by the service.
type Config struct {
	// Application configuration.
	ServerPort string
	AppName    string
	AppEnv     string

	// PostgreSQL configuration.
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string

	// Kafka configuration.
	//
	// DEVELOPMENT ONLY:
	// These values point to the local Docker Kafka instance.
	//
	// FUTURE REFACTOR:
	// Kafka broker configuration should eventually come from
	// the central TalentIQ Kafka/infrastructure configuration.
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
}

// Load reads configuration from environment variables.
//
// If an environment variable is not present, a safe
// local-development default value is used.
func Load() Config {
	return Config{
		// Application settings.
		ServerPort: getEnv(
			"SERVER_PORT",
			"8080",
		),
		AppName: getEnv(
			"APP_NAME",
			"talent-profile-intelligence",
		),
		AppEnv: getEnv(
			"APP_ENV",
			"development",
		),

		// PostgreSQL settings.
		DatabaseHost: getEnv(
			"DB_HOST",
			"localhost",
		),
		DatabasePort: getEnv(
			"DB_PORT",
			"5432",
		),
		DatabaseUser: getEnv(
			"DB_USER",
			"postgres",
		),
		DatabasePassword: getEnv(
			"DB_PASSWORD",
			"postgres",
		),
		DatabaseName: getEnv(
			"DB_NAME",
			"talent_db",
		),

		// Local Docker Kafka configuration.
		//
		// DEVELOPMENT ONLY:
		// Kafka currently runs locally on localhost:9092.
		//
		// FUTURE REFACTOR:
		// These values should eventually come from the
		// centralized TalentIQ Kafka configuration.
		KafkaBrokers: []string{
			getEnv(
				"KAFKA_BROKERS",
				"localhost:9092",
			),
		},

		KafkaTopic: getEnv(
			"KAFKA_TOPIC",
			"talent-profile-events",
		),

		KafkaGroupID: getEnv(
			"KAFKA_GROUP_ID",
			"talent-profile-intelligence",
		),
	}
}

// getEnv returns the environment variable value.
//
// If the variable does not exist or is empty,
// the provided default value is returned.
func getEnv(
	key string,
	defaultValue string,
) string {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}