package api

import (
	Handler "glamapp/src/handlers"

	"glamapp/src/database"

	"github.com/gofiber/fiber/v2"
)

func registerProfileRoutes(router fiber.Router, db *database.MongoDB) {
	profiles := router.Group("/profiles")
	profiles.Get("/", Handler.GetUsers(db))
	profiles.Get("/:id", Handler.GetUser(db))
	profiles.Post("/", Handler.CreateUser(db))
	profiles.Put("/:id", Handler.UpdateUser(db))
}
