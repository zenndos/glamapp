package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Avatar    string             `bson:"avatar"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type Post struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Content    string             `bson:"content"`
	AuthorID   primitive.ObjectID `bson:"author_id"`
	LikesCount int                `bson:"likes_count"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}
