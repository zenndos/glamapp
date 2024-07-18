package database

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client        *mongo.Client
	database      string
	users         *mongo.Collection
	notifications *mongo.Collection
	sessions      *mongo.Collection
	history       *mongo.Collection
	posts         *mongo.Collection

	logger zerolog.Logger
}

func NewMongoDB(uri, database string, logger zerolog.Logger) *MongoDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	db := client.Database(database)
	return &MongoDB{
		client:        client,
		database:      database,
		users:         db.Collection("users"),
		posts:         db.Collection("posts"),
		sessions:      db.Collection("sessions"),
		history:       db.Collection("history"),
		notifications: db.Collection("notifications"),

		logger: logger,
	}
}
