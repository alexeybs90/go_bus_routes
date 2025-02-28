package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alexeybs90/go_bus_routes/internal/route"
	"github.com/alexeybs90/go_bus_routes/pkg/logger"
)

type handlers struct {
	repository route.Repository
	logger     logger.Logger
}

func NewHandler(repository route.Repository, logger logger.Logger) *handlers {
	return &handlers{
		repository: repository,
		logger:     logger,
	}
}

// func (h *handlers) Register(router *httprouter.Router) {
// 	router.HandlerFunc(http.MethodGet, authorsURL, apperror.Middleware(h.GetList))
// }

func (h *handlers) GetList(w http.ResponseWriter, r *http.Request) {
	all, err := h.repository.FindAll(context.Background())
	if err != nil {
		w.WriteHeader(400)
		return
	}

	allBytes, err := json.Marshal(all)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

	return
}
