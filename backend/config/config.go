package config

import "os"

type Config struct {
	ServerPort  string
	Environment string
}

func Load() Config {
	port := os.Getenv("TALENTIQ_SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	environment := os.Getenv("TALENTIQ_ENV")
	if environment == "" {
		environment = "development"
	}

	return Config{
		ServerPort:  port,
		Environment: environment,
	}
}
