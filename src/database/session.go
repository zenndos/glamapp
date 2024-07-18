package database

import (
	"context"
	"time"

	"glamapp/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *MongoDB) CreateSession(userID primitive.ObjectID, token string, expiresAt time.Time) (*models.Session, error) {
	session := &models.Session{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		Token:        token,
		CreatedAt:    time.Now(),
		ExpiresAt:    expiresAt,
		LastActivity: time.Now(),
	}

	_, err := m.sessions.InsertOne(context.Background(), session)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to create session")
		return nil, err
	}

	return session, nil
}

func (m *MongoDB) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	err := m.sessions.FindOne(context.Background(), bson.M{"token": token}).Decode(&session)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to get session by token")
		return nil, err
	}
	return &session, nil
}

func (m *MongoDB) UpdateSessionActivity(sessionID primitive.ObjectID) error {
	_, err := m.sessions.UpdateOne(
		context.Background(),
		bson.M{"_id": sessionID},
		bson.M{"$set": bson.M{"last_activity": time.Now()}},
	)
	return err
}

func (m *MongoDB) DeleteSession(sessionID primitive.ObjectID) error {
	_, err := m.sessions.DeleteOne(context.Background(), bson.M{"_id": sessionID})
	return err
}

func (m *MongoDB) LogHistory(userID primitive.ObjectID) error {
	history := &models.History{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Timestamp: time.Now(),
	}

	_, err := m.history.InsertOne(context.Background(), history)
	return err
}
