package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/alexeybs90/go_bus_routes/internal/config"
	"github.com/alexeybs90/go_bus_routes/internal/handlers"
	"github.com/alexeybs90/go_bus_routes/internal/repository"
	"github.com/alexeybs90/go_bus_routes/pkg/logger"
	"github.com/alexeybs90/go_bus_routes/pkg/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	db     *pgxpool.Pool
	logger logger.Logger
	cfg    config.Config
	server *http.Server
}

func New(ctx context.Context, cfg config.Config) *App {
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
	handler.Register(router)

	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &App{
		logger: log,
		db:     client,
		cfg:    cfg,
		server: server,
	}
}

func (app *App) Run() error {
	app.logger.Info("starting server", slog.String("address", app.cfg.Server.Address))
	err := app.server.ListenAndServe()
	if err != nil {
		app.logger.Error(err.Error())
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
