package config

import "os"

type Config struct {
	ServerPort string
	AppName    string
	AppEnv     string
}

func Load() Config {
	return Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		AppName:    getEnv(
			"APP_NAME",
			"talent-profile-intelligence",
		),
		AppEnv: getEnv("APP_ENV", "development"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}