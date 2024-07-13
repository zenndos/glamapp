package handlers

import (
	"github.com/gofiber/fiber/v2"

	"glamapp/src/database"
	"glamapp/src/models"
)

func GetUser(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		profile, err := db.GetUser(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"detail": "Profile not found",
			})
		}

		return c.JSON(fiber.Map{
			"data": profile,
		})
	}
}

func GetUsers(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		profiles, err := db.GetUsers()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": "Internal Server Error",
			})
		}

		return c.JSON(fiber.Map{
			"data":  profiles,
			"count": len(profiles),
		})
	}
}

func CreateUser(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		profile := new(models.User)
		if err := c.BodyParser(profile); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "Invalid request body",
			})
		}

		id, err := db.CreateUser(profile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": "Failed to create profile",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id": id,
		})
	}
}

func UpdateUser(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		profile := new(models.User)
		if err := c.BodyParser(profile); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "Invalid request body",
			})
		}

		if err := db.UpdateUser(id, profile); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": "Failed to update profile",
			})
		}

		return c.JSON(fiber.Map{
			"detail": "Profile updated successfully",
		})
	}
}
