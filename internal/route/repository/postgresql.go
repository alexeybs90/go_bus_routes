package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexeybs90/go_bus_routes/internal/route"
	"github.com/alexeybs90/go_bus_routes/pkg/client/postgresql"
	"github.com/alexeybs90/go_bus_routes/pkg/logger"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	client postgresql.Storage
	logger logger.Logger
}

// Create implements route.Repository.
func (r *repository) Create(ctx context.Context, item *route.Route) error {
	sql := "INSERT INTO route (name) VALUES ($1) RETURNING id"
	if err := r.client.QueryRow(ctx, sql, item.Name).Scan(&item.Id); err != nil {
		r.logDB(err)
		return err
	}
	return nil
}

// Delete implements route.Repository.
func (r *repository) Delete(ctx context.Context, id int) error {
	sql := "DELETE FROM route WHERE id=$1"
	_, err := r.client.Query(ctx, sql, id)
	if err != nil {
		r.logDB(err)
		return err
	}
	return nil
}

// FindAll implements route.Repository.
func (r *repository) FindAll(ctx context.Context) ([]route.Route, error) {
	sql := "SELECT id, name FROM route ORDER BY name"
	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		r.logDB(err)
		return nil, err
	}
	defer rows.Close()

	items := make([]route.Route, 0)

	for rows.Next() {
		var item route.Route
		err = rows.Scan(&item.Id, &item.Name)
		if err != nil {
			r.logDB(err)
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// FindOne implements route.Repository.
func (r *repository) FindOne(ctx context.Context, id int) (route.Route, error) {
	var item route.Route
	sql := "SELECT id, name FROM route WHERE id=$1"
	if err := r.client.QueryRow(ctx, sql, id).Scan(&item.Id, &item.Name); err != nil {
		r.logDB(err)
		return item, err
	}
	return item, nil
}

// Update implements route.Repository.
func (r *repository) Update(ctx context.Context, item route.Route) error {
	sql := "UPDATE route SET name=$1 WHERE id=$2"
	_, err := r.client.Query(ctx, sql, item.Name, item.Id)
	if err != nil {
		r.logDB(err)
		return err
	}
	return nil
}

func NewRepository(client postgresql.Storage, logger logger.Logger) route.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) logDB(err error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		pgErr = err.(*pgconn.PgError)
		r.logger.Error(
			fmt.Sprintf(
				"SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState(),
			),
		)
	}
}
