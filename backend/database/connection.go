package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func NewConnection(ctx context.Context) (*pgxpool.Pool, error) {

	// Load .env file if it exists.
	_ = godotenv.Load()

	// Read DATABASE_URL.
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	// Create PostgreSQL connection pool.
	pool, err := pgxpool.New(ctx, databaseURL)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create database connection: %w",
			err,
		)
	}

	// Check database connection.
	err = pool.Ping(ctx)

	if err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"failed to connect to database: %w",
			err,
		)
	}

	return pool, nil
}
