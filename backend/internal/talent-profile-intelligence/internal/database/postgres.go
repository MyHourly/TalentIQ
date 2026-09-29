package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"talentiq/talent-profile-intelligence/internal/config"
)

// NewPostgres creates a PostgreSQL connection pool.
//
// The pool keeps database connections available for reuse,
// which is more efficient than opening a new connection
// for every database operation.
func NewPostgres(
	ctx context.Context,
	cfg config.Config,
) (*pgxpool.Pool, error) {

	// Build the PostgreSQL connection string.
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)

	// Create the PostgreSQL connection pool.
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create PostgreSQL connection pool: %w",
			err,
		)
	}

	// Use a timeout so application startup does not wait
	// indefinitely if PostgreSQL is unavailable.
	pingCtx, cancel := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer cancel()

	// Verify that PostgreSQL is reachable.
	if err := pool.Ping(pingCtx); err != nil {
		// Close the pool if the connection check fails.
		pool.Close()

		return nil, fmt.Errorf(
			"failed to connect to PostgreSQL: %w",
			err,
		)
	}

	return pool, nil
}
