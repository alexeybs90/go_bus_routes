package main

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/alexeybs90/go_bus_routes/internal/config"
	"github.com/alexeybs90/go_bus_routes/internal/route"
	"github.com/alexeybs90/go_bus_routes/internal/route/repository"
	"github.com/alexeybs90/go_bus_routes/pkg/client/postgresql"
	"github.com/alexeybs90/go_bus_routes/pkg/logger"
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

	client, err := postgresql.NewClient(context.Background(), cfg.Storage)
	if err != nil {
		log.Error(err.Error())
	}

	repRoute := repository.NewRepository(*client, log)

	item := route.Route{Id: 1, Name: "Автобус № 11 Купчино-Невский"}
	err = repRoute.Update(context.Background(), item)
	if err != nil {
		log.Error(err.Error())
	}
	items, err := repRoute.FindAll(context.Background())
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Println(items)
}
