package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"glamapp/src/api"
	Config "glamapp/src/config"
	"glamapp/src/database"
)

type App struct {
	*fiber.App

	DB     *database.MongoDB
	Logger zerolog.Logger
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logger := log.With().Str("service", "glamapp").Logger()

	config := Config.ReadConfig()

	db := database.NewMongoDB(config.DatabaseURI, config.Database)

	app := App{
		App:    fiber.New(),
		DB:     db,
		Logger: logger,
	}
	api_core_group := app.Group("/api")
	api_v1_group := api_core_group.Group("/v1")
	api.RegisterRoutes(api_v1_group, db)

	app.Listen(":" + config.AppPort)
}
