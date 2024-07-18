package main

import (
	"fmt"
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

	db := database.NewMongoDB(config.DatabaseURI, config.Database, logger)

	app := App{
		App:    fiber.New(),
		DB:     db,
		Logger: logger,
	}

	api.RegisterAuthRoutes(app.App, db, logger)

	apiV1 := app.Group("/api/v1")

	apiV1.Use(api.JWTMiddleware(db, logger))

	apiV1.Use(api.SessionMiddleware(db, logger))

	api.RegisterProfileRoutes(apiV1, db, logger)
	api.RegisterPostRoutes(apiV1, db, logger)

	address := fmt.Sprintf("%s:%s", config.AppHost, config.AppPort)
	logger.Info().Msgf("Starting server on %s", address)
	if err := app.Listen(address); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
