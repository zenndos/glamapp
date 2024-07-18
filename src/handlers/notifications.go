package handlers

import (
	"glamapp/src/database"
	"glamapp/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type NotificationHandler struct {
	DB     *database.MongoDB
	Logger zerolog.Logger
}

func NewNotificationHandler(db *database.MongoDB, logger zerolog.Logger) *NotificationHandler {
	return &NotificationHandler{
		DB:     db,
		Logger: logger,
	}
}

func (h *NotificationHandler) ReadAllNotifications(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	notifications, err := h.DB.GetAndDeleteUserNotifications(user.ID)
	if err != nil {
		h.Logger.Error().Err(err).Str("user_id", user.ID.Hex()).Msg("Failed to fetch and delete notifications")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process notifications")
	}

	return c.JSON(fiber.Map{
		"notifications": notifications,
		"count":         len(notifications),
	})
}
