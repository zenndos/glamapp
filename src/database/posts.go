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

func (m *MongoDB) CreatePost(post *models.Post, userID primitive.ObjectID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.LikesCount = 0

	// Insert the post
	result, err := m.posts.InsertOne(ctx, post)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to insert post")
		return "", fmt.Errorf("failed to insert post: %w", err)
	}

	postID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		m.logger.Error().Msg("Failed to get inserted ID")
		return "", fmt.Errorf("failed to get inserted ID")
	}

	// Update the user's posts array
	updateBody := bson.M{
		"$push": bson.M{"posts": postID.Hex()},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err = m.users.UpdateOne(ctx, bson.M{"_id": userID}, updateBody)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to update user's posts")
		return "", fmt.Errorf("failed to update user's posts: %w", err)
	}

	return postID.Hex(), nil
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

func (m *MongoDB) GetPosts() ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.posts.Find(ctx, bson.M{})
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

func (m *MongoDB) LikePost(postID string, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		m.logger.Error().Err(err).Str("post_id", postID).Str("user_id", userID.Hex()).Msg("Invalid post ID")
		return fmt.Errorf("invalid post ID: %w", err)
	}

	postUpdate := bson.M{
		"$inc":      bson.M{"likes_count": 1},
		"$addToSet": bson.M{"liked_by": userID},
		"$set":      bson.M{"updated_at": time.Now()},
	}
	postResult, err := m.posts.UpdateOne(
		ctx,
		bson.M{"_id": postObjectID, "liked_by": bson.M{"$ne": userID}},
		postUpdate,
	)
	if err != nil {
		m.logger.Error().Err(err).Str("post_id", postID).Str("user_id", userID.Hex()).Msg("Failed to update post")
		return fmt.Errorf("failed to update post: %w", err)
	}
	if postResult.ModifiedCount == 0 {
		return fmt.Errorf("post already liked or not found")
	}

	userUpdate := bson.M{
		"$addToSet": bson.M{"liked_posts": postID},
	}
	userResult, err := m.users.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		userUpdate,
	)
	if err != nil {
		m.logger.Error().Err(err).Str("post_id", postID).Str("user_id", userID.Hex()).Msg("Failed to update user")
		return fmt.Errorf("failed to update user: %w", err)
	}
	if userResult.ModifiedCount == 0 {
		m.logger.Warn().Str("post_id", postID).Str("user_id", userID.Hex()).Msg("User not found or post already in liked_posts")
	}

	m.logger.Info().Str("post_id", postID).Str("user_id", userID.Hex()).Msg("Post liked successfully")
	return nil
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

func (m *MongoDB) DeleteOne(collection *mongo.Collection, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to convert ID to ObjectID")
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
