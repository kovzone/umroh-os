package main

import (
	"context"
	"fmt"
	"time"

	"booking-svc/util/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// createPostgresPool creates a new postgres pool.
func createPostgresPool(postgresConfig config.Postgres) (*pgxpool.Pool, error) {
	pgConfig, err := pgxpool.ParseConfig(postgresConfig.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	pgConfig.MaxConns = int32(postgresConfig.Pool.MaxConns)
	pgConfig.MinConns = int32(postgresConfig.Pool.MinConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping postgres pool: %w", err)
	}

	return pool, nil
}

// createPostgresPoolWithRetry creates a new postgres pool with retry.
func createPostgresPoolWithRetry(postgresConfig config.Postgres) (*pgxpool.Pool, error) {
	maxRetries := postgresConfig.Pool.RetryMaxAttempts
	baseDelay := postgresConfig.Pool.RetryBaseDelay

	for attempt := 1; attempt <= maxRetries; attempt++ {
		pool, err := createPostgresPool(postgresConfig)
		if err == nil {
			return pool, nil
		}

		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt) * baseDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to postgres after %d attempts", maxRetries)
}
