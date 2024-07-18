package api

import (
	"glamapp/src/config"
	"glamapp/src/database"
	"glamapp/src/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func RegisterAuthRoutes(app *fiber.App, db *database.MongoDB, logger zerolog.Logger) {
	conf := config.ReadConfig()
	authHandler := handlers.NewAuthHandler(db, logger, conf.JWTSecret)

	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
}

func RegisterProfileRoutes(router fiber.Router, db *database.MongoDB, logger zerolog.Logger) {
	userHandler := handlers.NewUserHandler(db, logger)

	users := router.Group("/users")
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Patch("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
	users.Get("/:id/avatar", userHandler.GetAvatar)
}

func RegisterPostRoutes(router fiber.Router, db *database.MongoDB, logger zerolog.Logger) {
	postHandler := handlers.NewPostHandler(db, logger)

	posts := router.Group("/posts")
	posts.Post("/", postHandler.CreatePost)
	posts.Get("/", postHandler.GetPosts)
	posts.Get("/:id", postHandler.GetPost)
	posts.Post("/:id/like", postHandler.LikePost)
	posts.Delete("/:id", postHandler.DeletePost)
	posts.Get("/profile/:profileId", postHandler.GetPostsByProfile)
}

func RegisterNotificationRoutes(router fiber.Router, db *database.MongoDB, logger zerolog.Logger) {
	notificationHandler := handlers.NewNotificationHandler(db, logger)

	notifications := router.Group("/notifications")
	notifications.Get("/", notificationHandler.ReadAllNotifications)
}
