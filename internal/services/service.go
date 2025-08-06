package services

import (
	"context"

	"github.com/alexeybs90/go_bus_routes/internal/model"
	"github.com/alexeybs90/go_bus_routes/pkg/logger"
)

type busService struct {
	repository model.Repository
	logger     logger.Logger
}

func New(rep model.Repository, log logger.Logger) *busService {
	return &busService{
		repository: rep,
		logger:     log,
	}
}

func (s *busService) FindBus(ctx context.Context, fromId int, toId int) ([]model.Model, error) {
	return s.repository.FindBus(ctx, fromId, toId)
}
