package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alexeybs90/go_bus_routes/internal/route"
	"github.com/alexeybs90/go_bus_routes/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type handlers struct {
	repository route.Repository
	logger     logger.Logger
}

type response struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type responseRoute struct {
	response
	Route route.Route `json:"route,omitempty"`
}

type responseRoutes struct {
	response
	Routes []route.Route `json:"routes,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func NewHandler(repository route.Repository, logger logger.Logger) *handlers {
	return &handlers{
		repository: repository,
		logger:     logger,
	}
}

func (h *handlers) Register(router *chi.Mux) {
	router.Get("/api/routes", h.GetList)
	router.Get("/api/routes/{id}", h.FindOne)
	router.Post("/api/routes", h.Create)
	router.Put("/api/routes", h.Update)
	router.Delete("/api/routes/{id}", h.Delete)
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

func (h *handlers) GetList(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(
		slog.String("api", "handlers.GetList"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resp := json.NewEncoder(w)

	all, err := h.repository.FindAll(r.Context())
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	log.Info("done ok!")
	w.WriteHeader(http.StatusOK)
	resp.Encode(responseRoutes{
		response: response{Status: StatusOK},
		Routes:   all,
	})
}

func (h *handlers) FindOne(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(
		slog.String("api", "handlers.Create"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//id := r.URL.Query().Get("id")
	id := chi.URLParam(r, "id")
	if id == "" {
		err := errors.New("request param id is required")
		h.doServerError(log, err, w)
		return
	}
	route_id, err := strconv.Atoi(id)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	route, err := h.repository.FindOne(r.Context(), route_id)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	log.Info("done ok!")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseRoute{
		response: response{Status: StatusOK},
		Route:    route,
	})
}

func (h *handlers) Create(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(
		slog.String("api", "handlers.Create"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	route := route.Route{}
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &route); err != nil {
		h.doServerError(log, err, w)
		return
	}

	err = h.repository.Create(r.Context(), &route)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	log.Info("done ok!")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseRoute{
		response: response{Status: StatusOK},
		Route:    route,
	})
}

func (h *handlers) Update(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(
		slog.String("api", "handlers.Update"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	route := route.Route{}
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &route); err != nil {
		h.doServerError(log, err, w)
		return
	}

	err = h.repository.Update(r.Context(), route)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	log.Info("done ok!")
	w.WriteHeader(http.StatusOK)
	h.responseOK(w)
}

func (h *handlers) Delete(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(
		slog.String("api", "handlers.Delete"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := chi.URLParam(r, "id")
	if id == "" {
		err := errors.New("request param id is required")
		h.doServerError(log, err, w)
		return
	}
	route_id, err := strconv.Atoi(id)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	err = h.repository.Delete(r.Context(), route_id)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	log.Info("done ok!")
	w.WriteHeader(http.StatusOK)
	h.responseOK(w)
}
