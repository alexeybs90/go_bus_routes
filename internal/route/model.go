package route

import "github.com/alexeybs90/go_bus_routes/internal/station"

type Route struct {
	Id       int               `json:"id"`
	Name     string            `json:"name"`
	Stations []station.Station `json:"stations"`
}
