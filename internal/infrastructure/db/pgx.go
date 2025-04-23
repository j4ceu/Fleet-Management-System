package db

import (
	"context"
	"time"

	"internal/infrastructure/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return pgxpool.New(ctx, cfg.DbURL)
}
