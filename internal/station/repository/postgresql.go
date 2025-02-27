package repository

import (
	"context"

	"github.com/alexeybs90/go_bus_routes/internal/station"
	"github.com/alexeybs90/go_bus_routes/pkg/logger"
	"github.com/alexeybs90/go_bus_routes/pkg/storage/postgresql"
)

type repository struct {
	client *postgresql.Storage
	logger logger.Logger
}

// Create implements station.Repository.
func (r *repository) Create(ctx context.Context, item *station.Station) error {
	sql := "INSERT INTO station (name) VALUES ($1) RETURNING id"
	if err := r.client.QueryRow(ctx, sql, item.Name).Scan(&item.Id); err != nil {
		r.client.LogDB(err)
		return err
	}
	return nil
}

// Delete implements station.Repository.
func (r *repository) Delete(ctx context.Context, id int) error {
	sql := "DELETE FROM station WHERE id=$1"
	_, err := r.client.Query(ctx, sql, id)
	if err != nil {
		r.client.LogDB(err)
		return err
	}
	return nil
}

// FindAll implements station.Repository.
func (r *repository) FindAll(ctx context.Context) ([]station.Station, error) {
	sql := "SELECT id, name FROM station ORDER BY name"
	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		r.client.LogDB(err)
		return nil, err
	}
	defer rows.Close()

	items := make([]station.Station, 0)

	for rows.Next() {
		var item station.Station
		err = rows.Scan(&item.Id, &item.Name)
		if err != nil {
			r.client.LogDB(err)
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// FindOne implements station.Repository.
func (r *repository) FindOne(ctx context.Context, id int) (station.Station, error) {
	var item station.Station
	sql := "SELECT id, name FROM station WHERE id=$1"
	if err := r.client.QueryRow(ctx, sql, id).Scan(&item.Id, &item.Name); err != nil {
		r.client.LogDB(err)
		return item, err
	}
	return item, nil
}

// Update implements station.Repository.
func (r *repository) Update(ctx context.Context, item station.Station) error {
	sql := "UPDATE station SET name=$1 WHERE id=$2"
	_, err := r.client.Query(ctx, sql, item.Name, item.Id)
	if err != nil {
		r.client.LogDB(err)
		return err
	}
	return nil
}

func NewRepository(client *postgresql.Storage, logger logger.Logger) station.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
