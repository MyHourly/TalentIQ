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
}

// Load reads configuration from environment variables.
//
// If an environment variable is not present, a safe local-development
// default value is used.
func Load() Config {
	return Config{
		// Application settings.
		ServerPort: getEnv("SERVER_PORT", "8080"),
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
	}
}

// getEnv returns the environment variable value.
//
// If the variable does not exist or is empty, the provided
// default value is returned.
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
