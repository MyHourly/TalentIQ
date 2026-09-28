package database

import (
	"context"
	"fmt"
	"time"

	"talentiq/talent-profile-intelligence/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgres creates a PostgreSQL connection pool.
//
// A connection pool allows the application to reuse database
// connections instead of creating a new connection for every request.
func NewPostgres(
	ctx context.Context,
	cfg config.Config,
) (*pgxpool.Pool, error) {

	// Build the PostgreSQL connection string from application configuration.
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)

	// Create the connection pool.
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create PostgreSQL connection pool: %w",
			err,
		)
	}

	// Give PostgreSQL a limited amount of time to respond.
	pingCtx, cancel := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer cancel()

	// Verify that the database is actually reachable.
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"failed to connect to PostgreSQL: %w",
			err,
		)
	}

	return pool, nil
}