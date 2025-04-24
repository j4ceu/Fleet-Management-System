package location

import (
	"FleetManagementSystem/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SaveLocation(ctx context.Context, loc entity.Location) error
	GetLastLocation(ctx context.Context, vehicleID string) (entity.Location, error)
	GetLocationHistory(ctx context.Context, payload entity.LocationHistoryPayload) ([]entity.Location, error)
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

func (r *repository) GetLastLocation(ctx context.Context, vehicleID string) (entity.Location, error) {
	row := r.db.QueryRow(ctx, `
		SELECT vehicle_id, latitude, longitude, timestamp
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY timestamp DESC
		LIMIT 1`, vehicleID)

	var loc entity.Location
	err := row.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp)
	return loc, err
}

func (r *repository) GetLocationHistory(ctx context.Context, payload entity.LocationHistoryPayload) ([]entity.Location, error) {
	rows, err := r.db.Query(ctx, `
		SELECT vehicle_id, latitude, longitude, timestamp
		FROM vehicle_locations
		WHERE vehicle_id = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp ASC`, payload.VehicleID, payload.StartTime, payload.EndTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []entity.Location
	for rows.Next() {
		var loc entity.Location
		if err := rows.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil

}
