package repository

import (
	"context"

	"github.com/alexeybs90/go_bus_routes/internal/route"
	"github.com/alexeybs90/go_bus_routes/internal/station"
	"github.com/alexeybs90/go_bus_routes/pkg/logger"
	"github.com/alexeybs90/go_bus_routes/pkg/storage/postgresql"
)

type repository struct {
	client *postgresql.Storage
	logger logger.Logger
}

// Create implements route.Repository.
func (r *repository) Create(ctx context.Context, item *route.Route) error {
	sql := "INSERT INTO route (name) VALUES ($1) RETURNING id"
	if err := r.client.QueryRow(ctx, sql, item.Name).Scan(&item.Id); err != nil {
		r.LogDB(err)
		return err
	}
	return nil
}

// Delete implements route.Repository.
func (r *repository) Delete(ctx context.Context, id int) error {
	sql := "DELETE FROM route WHERE id=$1"
	_, err := r.client.Query(ctx, sql, id)
	if err != nil {
		r.LogDB(err)
		return err
	}
	return nil
}

// FindAll implements route.Repository.
func (r *repository) FindAll(ctx context.Context) ([]route.Route, error) {
	sqlStation := `SELECT s.id, s.name, rs.route_id FROM route_stations rs
			JOIN station s ON s.id=rs.station_id
	ORDER BY pos`
	rowsStation, err := r.client.Query(ctx, sqlStation)
	if err != nil {
		r.LogDB(err)
		return nil, err
	}
	defer rowsStation.Close()
	stationsByRouteId := make(map[int][]station.Station, 0)
	for rowsStation.Next() {
		var st station.Station
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

	routes := make([]route.Route, 0)

	for rows.Next() {
		var rt route.Route
		err = rows.Scan(&rt.Id, &rt.Name)
		if err != nil {
			r.LogDB(err)
			return nil, err
		}
		rt.Stations = []station.Station{}
		if arr, ok := stationsByRouteId[rt.Id]; ok {
			rt.Stations = arr
		}
		routes = append(routes, rt)
	}

	return routes, nil
}

// FindOne implements route.Repository.
func (r *repository) FindOne(ctx context.Context, id int) (route.Route, error) {
	var item route.Route
	sql := "SELECT id, name FROM route WHERE id=$1"
	if err := r.client.QueryRow(ctx, sql, id).Scan(&item.Id, &item.Name); err != nil {
		// r.client.LogDB(err)
		r.logger.Error(err.Error())
		return item, err
	}

	sqlStation := `SELECT s.id, s.name FROM route_stations rs
		JOIN station s ON s.id=rs.station_id WHERE rs.route_id=$1
		ORDER BY pos`
	rowsStation, err := r.client.Query(ctx, sqlStation, id)
	if err != nil {
		r.LogDB(err)
		return item, err
	}
	defer rowsStation.Close()
	item.Stations = []station.Station{}
	for rowsStation.Next() {
		var st station.Station
		err = rowsStation.Scan(&st.Id, &st.Name)
		if err != nil {
			r.LogDB(err)
			return item, err
		}
		item.Stations = append(item.Stations, st)
	}

	return item, nil
}

// Update implements route.Repository.
func (r *repository) Update(ctx context.Context, item route.Route) error {
	sql := "UPDATE route SET name=$1 WHERE id=$2"
	_, err := r.client.Query(ctx, sql, item.Name, item.Id)
	if err != nil {
		r.LogDB(err)
		return err
	}
	return nil
}

func (r *repository) LogDB(err error) {
	r.logger.Error(r.client.ErrorDetails(err))
}

func NewRepository(client *postgresql.Storage, logger logger.Logger) route.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
