package database

import (
	"context"
	"time"

	"glamapp/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *MongoDB) CreateUser(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := m.users.InsertOne(context.Background(), user)
	return err
}

func (m *MongoDB) GetUserByName(name string) (*models.User, error) {
	var user models.User
	err := m.users.FindOne(context.Background(), bson.M{"name": name}).Decode(&user)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to get user by name")
		return nil, err
	}
	return &user, nil
}

func (m *MongoDB) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := m.users.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to get user by ID")
		return nil, err
	}

	return &user, nil
}

func (m *MongoDB) UpdateUser(id string, updateData *models.User, updatedFields map[string]bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateFields := bson.M{}
	if updatedFields["name"] {
		updateFields["name"] = updateData.Name
	}
	if updatedFields["password_hash"] {
		updateFields["password_hash"] = updateData.PasswordHash
	}

	if updatedFields["avatar_data"] {
		updateFields["avatar_data"] = updateData.AvatarData
	}
	if updatedFields["avatar_type"] {
		updateFields["avatar_type"] = updateData.AvatarType
	}
	updateFields["updated_at"] = time.Now()

	update := bson.M{
		"$set": updateFields,
	}

	_, err = m.users.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)

	return err
}

func (m *MongoDB) GetUser(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to convert ID to ObjectID")
		return nil, err
	}

	var user models.User
	err = m.users.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to get user by ID")
		return nil, err
	}

	return &user, nil
}

func (m *MongoDB) GetUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.users.Find(ctx, bson.M{})
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to fetch users")
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		m.logger.Error().Err(err).Msg("Failed to decode users")
		return nil, err
	}

	return users, nil
}

func (m *MongoDB) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to convert ID to ObjectID")
		return err
	}

	_, err = m.users.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to delete user")
		return err
	}

	return nil
}
