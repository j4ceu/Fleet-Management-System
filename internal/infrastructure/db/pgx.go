package db

import (
	"FleetManagementSystem/internal/infrastructure/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(cfg *config.Config) (*pgxpool.Pool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Close the pool if unable to ping the database
		return nil, fmt.Errorf("unable to reach database: %v", err)
	}

	return pool, nil
}
