package main

import (
	"context"
	"fmt"

	"github.com/alexeybs90/go_bus_routes/internal/app"
	"github.com/alexeybs90/go_bus_routes/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		panic(err)
	}

	app := app.New(context.Background(), cfg)
	err = app.Run()
	if err != nil {
		fmt.Println(err.Error())
	}

}
