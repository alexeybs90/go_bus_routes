package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alexeybs90/go_bus_routes/internal/model"
	"github.com/alexeybs90/go_bus_routes/pkg/logger"
	"github.com/go-chi/chi/v5"
)

const (
	contentType   = "application/json; charset=UTF-8"
	routeEntity   = "route"
	stationEntity = "station"
)

type handlers struct {
	repository model.Repository
	logger     logger.Logger
}

type response struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type responseItem struct {
	response
	Item model.Model `json:"item,omitempty"`
}

type responseItems struct {
	response
	Items []model.Model `json:"items,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func NewHandler(repository model.Repository, logger logger.Logger) *handlers {
	return &handlers{
		repository: repository,
		logger:     logger,
	}
}

func (h *handlers) Register(router *chi.Mux) {
	router.Get("/api/stations", h.GetStations)
	router.Get("/api/stations/{id}", h.GetStation)
	router.Post("/api/stations", h.CreateStation)
	router.Put("/api/stations", h.UpdateStation)
	router.Delete("/api/stations/{id}", h.DeleteStation)

	router.Get("/api/routes", h.GetRoutes)
	router.Get("/api/routes/{id}", h.GetRoute)
	router.Post("/api/routes", h.CreateRoute)
	router.Put("/api/routes", h.UpdateRoute)
	router.Delete("/api/routes/{id}", h.DeleteRoute)
}

func (h *handlers) responseError(err error, w http.ResponseWriter) {
	json.NewEncoder(w).Encode(response{
		Status: StatusError,
		Error:  err.Error(),
	})
}

func (h *handlers) responseOK(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(response{
		Status: StatusOK,
	})
}

func (h *handlers) doServerError(log logger.Logger, err error, w http.ResponseWriter) {
	log.Error(err.Error())
	w.WriteHeader(http.StatusBadGateway)
	h.responseError(err, w)
}

func (h *handlers) GetRoutes(w http.ResponseWriter, r *http.Request) {
	h.GetList(w, r, routeEntity)
}

func (h *handlers) GetStations(w http.ResponseWriter, r *http.Request) {
	h.GetList(w, r, stationEntity)
}

func (h *handlers) GetRoute(w http.ResponseWriter, r *http.Request) {
	h.GetOne(w, r, routeEntity)
}

func (h *handlers) GetStation(w http.ResponseWriter, r *http.Request) {
	h.GetOne(w, r, stationEntity)
}

func (h *handlers) CreateRoute(w http.ResponseWriter, r *http.Request) {
	route := &model.Route{}
	h.Create(w, r, route)
}

func (h *handlers) CreateStation(w http.ResponseWriter, r *http.Request) {
	station := &model.Station{}
	h.Create(w, r, station)
}

func (h *handlers) UpdateRoute(w http.ResponseWriter, r *http.Request) {
	item := &model.Route{}
	h.Update(w, r, item)
}

func (h *handlers) UpdateStation(w http.ResponseWriter, r *http.Request) {
	item := &model.Station{}
	h.Update(w, r, item)
}

func (h *handlers) DeleteRoute(w http.ResponseWriter, r *http.Request) {
	item := &model.Route{}
	h.Delete(w, r, item)
}

func (h *handlers) DeleteStation(w http.ResponseWriter, r *http.Request) {
	item := &model.Station{}
	h.Delete(w, r, item)
}
