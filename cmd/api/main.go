package main

import (
	"context"
	"fmt"
	"net/http"

	"log/slog"

	"github.com/alexeybs90/go_bus_routes/internal/config"
	"github.com/alexeybs90/go_bus_routes/internal/route/handlers"
	"github.com/alexeybs90/go_bus_routes/internal/route/repository"
	"github.com/alexeybs90/go_bus_routes/pkg/logger"
	"github.com/alexeybs90/go_bus_routes/pkg/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	log := logger.NewLogger(cfg.Env)
	log.Info("app started", slog.String("key", cfg.Env))
	log.Debug("debug messages are enabled")
	log.Error("error messages are enabled")

	client, err := postgresql.NewClient(context.Background(), cfg.Storage)
	if err != nil {
		log.Error(err.Error())
	}

	repRoute := repository.NewRepository(client, log)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	handler := handlers.NewHandler(repRoute, log)
	router.Get("/api/routes", handler.GetList)

	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	log.Info("starting server")
	err = server.ListenAndServe()
	if err != nil {
		log.Error(err.Error())
	}
}
