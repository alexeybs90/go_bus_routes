package model

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, item Model) error
	GetRoutes(ctx context.Context) ([]Model, error)
	GetRoute(ctx context.Context, id int) (Model, error)
	GetStation(ctx context.Context, id int) (Model, error)
	GetStations(ctx context.Context) ([]Model, error)
	Update(ctx context.Context, r Model) error
	Delete(ctx context.Context, r Model) error
	FindBus(ctx context.Context, fromId int, toId int) ([]Model, error)
}
