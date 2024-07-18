package database

import (
	"context"
	"fmt"
	"glamapp/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *MongoDB) GetAndDeleteUserNotifications(userID primitive.ObjectID) ([]models.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all notifications for the user
	cursor, err := m.notifications.Find(ctx, bson.M{"liked_by": userID})
	if err != nil {
		m.logger.Error().Err(err).Str("user_id", userID.Hex()).Msg("Failed to fetch notifications")
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}
	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		m.logger.Error().Err(err).Str("user_id", userID.Hex()).Msg("Failed to decode notifications")
		return nil, fmt.Errorf("failed to decode notifications: %w", err)
	}

	// Delete all notifications for the user
	_, err = m.notifications.DeleteMany(ctx, bson.M{"liked_by": userID})
	if err != nil {
		m.logger.Error().Err(err).Str("user_id", userID.Hex()).Msg("Failed to delete notifications")
		return notifications, fmt.Errorf("failed to delete notifications: %w", err)
	}

	m.logger.Info().Str("user_id", userID.Hex()).Int("count", len(notifications)).Msg("Notifications fetched and deleted")
	return notifications, nil
}
