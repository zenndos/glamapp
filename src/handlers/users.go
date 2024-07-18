package handlers

import (
	"github.com/gofiber/fiber/v2"

	"glamapp/src/database"
	"glamapp/src/models"

	"github.com/rs/zerolog"
)

type UserHandler struct {
	DB     *database.MongoDB
	Logger zerolog.Logger
}

func NewUserHandler(db *database.MongoDB, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		DB:     db,
		Logger: logger,
	}
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	h.Logger.Debug().Str("id", id).Msg("Updating user")

	user := new(models.User)
	updatedFields, err := user.Parse(c, true)
	if err != nil {
		h.Logger.Error().Err(err).Str("id", id).Msg("Failed to parse user update data")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.DB.UpdateUser(id, user, updatedFields); err != nil {
		h.Logger.Error().Err(err).Str("id", id).Msg("Failed to update user in database")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user")
	}

	h.Logger.Info().Str("id", id).Msg("User updated successfully")
	return c.JSON(fiber.Map{
		"detail": "User updated successfully",
	})
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.DB.GetUsers()
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to fetch users")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch users")
	}

	baseURL := c.BaseURL()
	var responseUsers []map[string]interface{}
	for _, user := range users {
		responseUsers = append(responseUsers, user.ToResponse(baseURL))
	}

	return c.JSON(fiber.Map{
		"users": responseUsers,
		"count": len(responseUsers),
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.DB.GetUser(id)
	if err != nil {
		h.Logger.Error().Err(err).Str("id", id).Msg("Failed to get user")
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	baseURL := c.BaseURL()
	response := user.ToResponse(baseURL)

	return c.JSON(response)
}

func (h *UserHandler) Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	return c.JSON(user.ToResponse(c.BaseURL()))
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	h.Logger.Debug().Str("id", id).Msg("Deleting user")

	err := h.DB.DeleteUser(id)
	if err != nil {
		h.Logger.Error().Err(err).Str("id", id).Msg("Failed to delete user from database")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete user")
	}

	h.Logger.Info().Str("id", id).Msg("User deleted successfully")
	return c.JSON(fiber.Map{
		"detail": "User deleted successfully",
	})
}

func (h *UserHandler) GetAvatar(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.DB.GetUser(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Profile not found")
	}

	if len(user.AvatarData) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Avatar not found")
	}

	c.Set("Content-Type", user.AvatarType)
	return c.Send(user.AvatarData)
}
