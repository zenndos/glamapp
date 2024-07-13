package main

import "github.com/gofiber/fiber/v2"

type App struct {
	*fiber.App
}

func main() {
	app := App{fiber.New()}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Listen(":3000")
}
