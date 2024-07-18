package api

import (
	"glamapp/src/config"
	"glamapp/src/database"
	"glamapp/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SessionMiddleware(db *database.MongoDB, logger zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)
		if !ok {
			logger.Debug().Msg("No user found in context, skipping session logging")
			return c.Next()
		}

		logger.Debug().Str("user_id", user.ID.Hex()).Msg("User found in context, logging history")

		err := db.LogHistory(user.ID)
		if err != nil {
			logger.Error().Err(err).Str("user_id", user.ID.Hex()).Msg("Failed to log history")
		} else {
			logger.Debug().Str("user_id", user.ID.Hex()).Msg("Successfully logged user history")
		}

		return c.Next()
	}
}

func JWTMiddleware(db *database.MongoDB, logger zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			logger.Warn().Str("path", c.Path()).Msg("Missing authorization token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization token"})
		}

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			conf := config.ReadConfig()
			return []byte(conf.JWTSecret), nil
		})

		if err != nil {
			logger.Error().Err(err).Str("token", token).Msg("Invalid or expired token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		userID, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
		if err != nil {
			logger.Error().Err(err).Str("user_id", claims["user_id"].(string)).Msg("Invalid user ID in token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
		}

		user, err := db.GetUserByID(userID)
		if err != nil {
			logger.Error().Err(err).Str("user_id", userID.Hex()).Msg("User not found")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		logger.Debug().Str("user_id", userID.Hex()).Msg("User authenticated successfully")
		c.Locals("user", user)

		return c.Next()
	}
}
