package handlers

import (
	"glamapp/src/database"

	"github.com/gofiber/fiber/v2"
)

func GetProfiles(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		profiles, err := db.GetProfiles()
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

func GetProfile(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		profile, err := db.GetProfile(id)
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

func CreateProfile(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		profile := new(database.Profile)
		if err := c.BodyParser(profile); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "Invalid request body",
			})
		}

		id, err := db.CreateProfile(profile)
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

func UpdateProfile(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		profile := new(database.Profile)
		if err := c.BodyParser(profile); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "Invalid request body",
			})
		}

		if err := db.UpdateProfile(id, profile); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": "Failed to update profile",
			})
		}

		return c.JSON(fiber.Map{
			"detail": "Profile updated successfully",
		})
	}
}
