package api

import (
	Handler "glamapp/src/handlers"

	"glamapp/src/database"

	"github.com/gofiber/fiber/v2"
)

func registerProfileRoutes(router fiber.Router, db *database.MongoDB) {
	profiles := router.Group("/profiles")
	profiles.Get("/", Handler.GetProfiles(db))
	profiles.Get("/:id", Handler.GetProfile(db))
	profiles.Post("/", Handler.CreateProfile(db))
	profiles.Put("/:id", Handler.UpdateProfile(db))
}
