package handlers

import (
	"glamapp/src/database"
	"glamapp/src/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	DB        *database.MongoDB
	Logger    zerolog.Logger
	JWTSecret string
}

func NewAuthHandler(db *database.MongoDB, logger zerolog.Logger, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		DB:        db,
		Logger:    logger,
		JWTSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := user.SetPassword(user.PasswordHash); err != nil {
		h.Logger.Error().Err(err).Msg("Failed to hash password")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	if err := h.DB.CreateUser(&user); err != nil {
		h.Logger.Error().Err(err).Msg("Failed to create user")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.DB.GetUserByUsername(loginData.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !user.CheckPassword(loginData.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to generate token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	return c.JSON(fiber.Map{"token": t})
}
