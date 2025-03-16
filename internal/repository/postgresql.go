package repository

import (
	"context"
	"fmt"

	"github.com/alexeybs90/go_bus_routes/internal/model"
	"github.com/alexeybs90/go_bus_routes/pkg/logger"
	"github.com/alexeybs90/go_bus_routes/pkg/storage/postgresql"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	client *pgxpool.Pool
	logger logger.Logger
}

func (r *repository) Create(ctx context.Context, item model.Model) error {
	sql := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", item.DBTable())
	var id int
	if err := r.client.QueryRow(ctx, sql, item.GetName()).Scan(&id); err != nil {
		r.LogDB(err)
		return err
	}
	item.SetID(id)
	return nil
}

func (r *repository) Delete(ctx context.Context, item model.Model) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE id=$1", item.DBTable())
	_, err := r.client.Query(ctx, sql, item.GetID())
	if err != nil {
		r.LogDB(err)
		return err
	}
	return nil
}

func (r *repository) GetRoutes(ctx context.Context) ([]model.Model, error) {
	sqlStation := `SELECT s.id, s.name, rs.route_id FROM route_stations rs
			JOIN station s ON s.id=rs.station_id
	ORDER BY pos`
	rowsStation, err := r.client.Query(ctx, sqlStation)
	if err != nil {
		r.LogDB(err)
		return nil, err
	}
	defer rowsStation.Close()
	stationsByRouteId := make(map[int][]model.Station, 0)
	for rowsStation.Next() {
		var st model.Station
		var routeId int
		err = rowsStation.Scan(&st.Id, &st.Name, &routeId)
		if err != nil {
			r.LogDB(err)
			return nil, err
		}
		stationsByRouteId[routeId] = append(stationsByRouteId[routeId], st)
	}

	sql := "SELECT id, name FROM route ORDER BY name"
	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		r.LogDB(err)
		return nil, err
	}
	defer rows.Close()

	routes := make([]model.Model, 0)

	for rows.Next() {
		var rt model.Route
		err = rows.Scan(&rt.Id, &rt.Name)
		if err != nil {
			r.LogDB(err)
			return nil, err
		}
		rt.Stations = []model.Station{}
		if arr, ok := stationsByRouteId[rt.Id]; ok {
			rt.Stations = arr
		}
		routes = append(routes, &rt)
	}

	return routes, nil
}

func (r *repository) GetRoute(ctx context.Context, id int) (model.Model, error) {
	var item model.Route
	sql := "SELECT id, name FROM route WHERE id=$1"
	if err := r.client.QueryRow(ctx, sql, id).Scan(&item.Id, &item.Name); err != nil {
		r.LogDB(err)
		return &item, err
	}

	sqlStation := `SELECT s.id, s.name FROM route_stations rs
		JOIN station s ON s.id=rs.station_id WHERE rs.route_id=$1
		ORDER BY pos`
	rowsStation, err := r.client.Query(ctx, sqlStation, id)
	if err != nil {
		r.LogDB(err)
		return &item, err
	}
	defer rowsStation.Close()
	item.Stations = []model.Station{}
	for rowsStation.Next() {
		var st model.Station
		err = rowsStation.Scan(&st.Id, &st.Name)
		if err != nil {
			r.LogDB(err)
			return &item, err
		}
		item.Stations = append(item.Stations, st)
	}

	return &item, nil
}

func (r *repository) GetStation(ctx context.Context, id int) (model.Model, error) {
	var item model.Station
	sql := "SELECT id, name FROM station WHERE id=$1"
	if err := r.client.QueryRow(ctx, sql, id).Scan(&item.Id, &item.Name); err != nil {
		r.LogDB(err)
		return &item, err
	}
	return &item, nil
}

func (r *repository) GetStations(ctx context.Context) ([]model.Model, error) {
	sql := "SELECT id, name FROM station ORDER BY name"
	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		r.LogDB(err)
		return nil, err
	}
	defer rows.Close()

	items := make([]model.Model, 0)

	for rows.Next() {
		var item model.Station
		err = rows.Scan(&item.Id, &item.Name)
		if err != nil {
			r.LogDB(err)
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

func (r *repository) Update(ctx context.Context, item model.Model) error {
	sql := fmt.Sprintf("UPDATE %s SET name=$1 WHERE id=$2", item.DBTable())
	_, err := r.client.Query(ctx, sql, item.GetName(), item.GetID())
	if err != nil {
		r.LogDB(err)
		return err
	}
	return nil
}

func (r *repository) LogDB(err error) {
	r.logger.Error(postgresql.ErrorDetails(err))
}

func NewRepository(client *pgxpool.Pool, logger logger.Logger) model.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
