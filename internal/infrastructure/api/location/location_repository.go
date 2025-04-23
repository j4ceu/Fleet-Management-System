package location

import (
	"FleetManagementSystem/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SaveLocation(ctx context.Context, loc entity.Location) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) SaveLocation(ctx context.Context, loc entity.Location) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp)
		VALUES ($1, $2, $3, $4)`, loc.VehicleID, loc.Latitude, loc.Longitude, loc.Timestamp)
	return err
}
