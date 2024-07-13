package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client     *mongo.Client
	database   string
	collection string
}

type Profile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Avatar    string             `bson:"avatar"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func NewMongoDB(uri, database, collection string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (m *MongoDB) getCollection() *mongo.Collection {
	return m.client.Database(m.database).Collection(m.collection)
}

func (m *MongoDB) GetProfile(id string) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var profile Profile
	err = m.getCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (m *MongoDB) GetProfiles() ([]Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.getCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var profiles []Profile
	if err = cursor.All(ctx, &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (m *MongoDB) CreateProfile(profile *Profile) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	result, err := m.getCollection().InsertOne(ctx, profile)
	if err != nil {
		return "", err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", mongo.ErrNoDocuments
	}

	return id.Hex(), nil
}

func (m *MongoDB) UpdateProfile(id string, updateData *Profile) error {
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

	_, err = m.getCollection().UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)

	return err
}
