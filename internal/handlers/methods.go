package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alexeybs90/go_bus_routes/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *handlers) Create(w http.ResponseWriter, r *http.Request, item model.Model) {
	log := h.logger.With(
		slog.String("api", "handlers.Create"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	w.Header().Set("Content-Type", contentType)

	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), item); err != nil {
		h.doServerError(log, err, w)
		return
	}

	err = h.repository.Create(r.Context(), item)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	log.Info("done ok!")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseItem{
		response: response{Status: StatusOK},
		Item:     item,
	})
}

func (h *handlers) Update(w http.ResponseWriter, r *http.Request, item model.Model) {
	log := h.logger.With(
		slog.String("api", "handlers.Update"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	w.Header().Set("Content-Type", contentType)

	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &item); err != nil {
		h.doServerError(log, err, w)
		return
	}

	err = h.repository.Update(r.Context(), item)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	log.Info("done ok!")
	w.WriteHeader(http.StatusOK)
	h.responseOK(w)
}

func (h *handlers) Delete(w http.ResponseWriter, r *http.Request, item model.Model) {
	log := h.logger.With(
		slog.String("api", "handlers.Delete"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	w.Header().Set("Content-Type", contentType)

	id := chi.URLParam(r, "id")
	if id == "" {
		err := errors.New("request param id is required")
		h.doServerError(log, err, w)
		return
	}
	itemId, err := strconv.Atoi(id)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	item.SetID(itemId)

	err = h.repository.Delete(r.Context(), item)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	log.Info("done ok!")
	w.WriteHeader(http.StatusOK)
	h.responseOK(w)
}

func (h *handlers) GetOne(w http.ResponseWriter, r *http.Request, entity string) {
	log := h.logger.With(
		slog.String("api", "handlers.GetOne"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	w.Header().Set("Content-Type", contentType)

	//id := r.URL.Query().Get("id")
	id := chi.URLParam(r, "id")
	if id == "" {
		err := errors.New("request param id is required")
		h.doServerError(log, err, w)
		return
	}
	itemId, err := strconv.Atoi(id)
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	var item model.Model

	switch entity {
	case routeEntity:
		item, err = h.repository.GetRoute(r.Context(), itemId)
	case stationEntity:
		item, err = h.repository.GetStation(r.Context(), itemId)
	default:
		h.doServerError(log, errors.New("wrong entity error"), w)
		return
	}
	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	log.Info("done ok!")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseItem{
		response: response{Status: StatusOK},
		Item:     item,
	})
}

func (h *handlers) GetList(w http.ResponseWriter, r *http.Request, entity string) {
	log := h.logger.With(
		slog.String("api", "handlers.GetList"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	w.Header().Set("Content-Type", contentType)

	resp := json.NewEncoder(w)

	var all []model.Model
	var err error

	switch entity {
	case routeEntity:
		all, err = h.repository.GetRoutes(r.Context())
	case stationEntity:
		all, err = h.repository.GetStations(r.Context())
	default:
		h.doServerError(log, errors.New("wrong entity error"), w)
		return
	}

	if err != nil {
		h.doServerError(log, err, w)
		return
	}

	log.Info("done ok!")
	w.WriteHeader(http.StatusOK)
	resp.Encode(responseItems{
		response: response{Status: StatusOK},
		Items:    all,
	})
}
