package database

import (
	"context"
	"fmt"
	"glamapp/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoDB) CreatePost(post *models.Post) (string, error) {
	ctx := context.Background()

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.LikesCount = 0

	result, err := m.posts.InsertOne(ctx, post)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to insert post")
		return "", fmt.Errorf("failed to insert post: %w", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		m.logger.Error().Msg("Failed to get inserted ID")
		return "", fmt.Errorf("failed to get inserted ID: %w", mongo.ErrNoDocuments)
	}

	update := bson.M{
		"$push": bson.M{"posts": id},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err = m.profiles.UpdateOne(ctx, bson.M{"_id": post.AuthorID}, update)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to update profile")
		// Note: At this point, the post has been created but the profile hasn't been updated.
		// You might want to implement some kind of cleanup or compensating action here.
		return "", fmt.Errorf("failed to update profile: %w", err)
	}

	return id.Hex(), nil
}

func (m *MongoDB) GetPost(id string) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var post models.Post
	err = m.posts.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (m *MongoDB) UpdatePost(id string, updateData *models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateData.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"content":     updateData.Content,
			"likes_count": updateData.LikesCount,
			"updated_at":  updateData.UpdatedAt,
		},
	}

	_, err = m.posts.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)

	return err
}

func (m *MongoDB) DeletePost(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.posts.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (m *MongoDB) GetPostsByProfile(profileID string) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return nil, err
	}

	cursor, err := m.posts.Find(ctx, bson.M{"author_id": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *MongoDB) DeleteProfile(id string) error {
	return m.DeleteOne(m.profiles, id)
}

func (m *MongoDB) DeleteOne(collection *mongo.Collection, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
