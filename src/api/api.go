package api

import (
	"glamapp/src/handlers"
	Handler "glamapp/src/handlers"

	"glamapp/src/database"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, db *database.MongoDB) {
	registerProfileRoutes(app, db)
	registerPostRoutes(app, db)
}

func registerProfileRoutes(router fiber.Router, db *database.MongoDB) {
	profiles := router.Group("/profiles")
	profiles.Get("/", Handler.GetProfiles(db))
	profiles.Get("/:id", Handler.GetProfile(db))
	profiles.Post("/", Handler.CreateProfile(db))
	profiles.Put("/:id", Handler.UpdateProfile(db))
}

func registerPostRoutes(router fiber.Router, db *database.MongoDB) {
	posts := router.Group("/posts")
	posts.Post("/", handlers.CreatePost(db))
	posts.Get("/:id", handlers.GetPost(db))
	posts.Put("/:id", handlers.UpdatePost(db))
	posts.Delete("/:id", handlers.DeletePost(db))
	posts.Get("/profile/:profileId", handlers.GetPostsByProfile(db))
}
