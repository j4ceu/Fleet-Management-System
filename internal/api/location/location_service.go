package location

import (
	"FleetManagementSystem/internal/entity"
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type Service interface {
	StoreLocation(ctx context.Context, loc entity.Location) error
	GetLatestLocation(ctx context.Context, vehicleID string) (entity.Location, error)
	GetLocationHistory(ctx context.Context, payload entity.LocationHistoryPayload) ([]entity.Location, error)
}

type service struct {
	repo Repository
	log  *logrus.Logger
}

func NewService(r Repository, log *logrus.Logger) Service {
	return &service{repo: r, log: log}
}

func (s *service) StoreLocation(ctx context.Context, loc entity.Location) error {
	funcName := "StoreLocation"

	s.log.Infof("[INFO][%s] starting", funcName)
	s.log.Infof("[INFO][%s] validate payload", funcName)
	if loc.VehicleID == "" || loc.Latitude == 0 || loc.Longitude == 0 || loc.Timestamp == 0 {
		s.log.Infof("[ERROR][%s] validate payload fail")
		return errors.New("invalid location data: all fields are required")
	}

	s.log.Infof("[INFO][%s] start SaveLocation to database", funcName)
	err := s.repo.SaveLocation(ctx, loc)
	if err != nil {
		s.log.Infof("[INFO][%s] SaveLocation error: %s", funcName, err.Error())
		return err
	}

	s.log.Infof("[INFO][%s] finished", funcName)

	return nil
}

func (s *service) GetLatestLocation(ctx context.Context, vehicleID string) (entity.Location, error) {
	funcName := "GetLatestLocation"

	s.log.Infof("[INFO][%s] starting", funcName)
	s.log.Infof("[INFO][%s] validate payload", funcName)
	if vehicleID == "" {
		s.log.Infof("[ERROR][%s] vehicle_id is required")
		return entity.Location{}, errors.New("vehicle_id is required")
	}

	s.log.Infof("[INFO][%s] start query from db", funcName)
	loc, err := s.repo.GetLastLocation(ctx, vehicleID)
	if err != nil {
		s.log.Infof("[ERROR][%s] Error GetLastLocation on db : %s", funcName, err.Error())
		return entity.Location{}, err
	}
	if loc.VehicleID == "" {
		s.log.Infof("[ERROR][%s] location not found", funcName)
		return entity.Location{}, errors.New("location not found")
	}

	s.log.Infof("[INFO][%s] finished", funcName)
	return loc, nil
}

func (s *service) GetLocationHistory(ctx context.Context, payload entity.LocationHistoryPayload) ([]entity.Location, error) {
	funcName := "GetLocationHistory"

	s.log.Infof("[INFO][%s] starting", funcName)
	s.log.Infof("[INFO][%s] validate payload", funcName)
	if payload.VehicleID == "" {
		s.log.Infof("[ERROR][%s] vehicle_id is requred", funcName)
		return nil, errors.New("vehicle_id is required")
	}
	if payload.EndTime < payload.StartTime {
		s.log.Infof("[ERROR][%s] invalid range: end must be after start", funcName)
		return nil, errors.New("invalid range: end must be after start")
	}

	location, err := s.repo.GetLocationHistory(ctx, payload)
	if err != nil {
		s.log.Infof("[ERROR][%s] Error GetLocationHistory on db : %s", funcName, err.Error())
		return nil, err
	}

	s.log.Infof("[INFO][%s] finished", funcName)
	return location, nil

}
