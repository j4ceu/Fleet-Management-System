package location

import (
	"FleetManagementSystem/internal/entity"
	"context"
)

type Service interface {
	StoreLocation(ctx context.Context, loc entity.Location) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) StoreLocation(ctx context.Context, loc entity.Location) error {

	return s.repo.SaveLocation(ctx, loc)
}
