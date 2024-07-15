package database

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"glamapp/src/models"
)

type MongoDB struct {
	client   *mongo.Client
	database string
	profiles *mongo.Collection
	posts    *mongo.Collection
}

func NewMongoDB(uri, database string) *MongoDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	db := client.Database(database)
	return &MongoDB{
		client:   client,
		database: database,
		profiles: db.Collection("profiles"),
		posts:    db.Collection("posts"),
	}
}

func (m *MongoDB) GetProfile(id string) (*models.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var profile models.Profile
	err = m.profiles.FindOne(ctx, bson.M{"_id": objectID}).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (m *MongoDB) GetProfiles() ([]models.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.profiles.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var profiles []models.Profile
	if err = cursor.All(ctx, &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (m *MongoDB) CreateProfile(profile *models.Profile) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	result, err := m.profiles.InsertOne(ctx, profile)
	if err != nil {
		return "", err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", mongo.ErrNoDocuments
	}

	return id.Hex(), nil
}

func (m *MongoDB) UpdateProfile(id string, updateData *models.Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateData.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":       updateData.Name,
			"avatar":     updateData.Avatar,
			"updated_at": updateData.UpdatedAt,
		},
	}

	_, err = m.profiles.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)

	return err
}

func (m *MongoDB) CreatePost(post *models.Post) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.LikesCount = 0 // Initialize likes count to 0

	result, err := m.posts.InsertOne(ctx, post)
	if err != nil {
		return "", err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", mongo.ErrNoDocuments
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
