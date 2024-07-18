package models

import (
	"fmt"
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	AvatarData   []byte             `bson:"avatar_data" json:"-"`
	AvatarType   string             `bson:"avatar_type" json:"-"`
	Posts        []string           `bson:"posts" json:"posts"`
	LikedPosts   []string           `bson:"liked_posts" json:"liked_posts"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewUser() *User {
	return &User{
		Posts:      []string{},
		LikedPosts: []string{},
	}
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) ToResponse(baseURL string) map[string]interface{} {
	return map[string]interface{}{
		"id":          u.ID,
		"name":        u.Name,
		"avatar":      baseURL + "/api/v1/users/" + u.ID.Hex() + "/avatar",
		"posts":       u.Posts,
		"liked_posts": u.LikedPosts,
		"created_at":  u.CreatedAt,
		"updated_at":  u.UpdatedAt,
	}
}

func (u *User) Parse(c *fiber.Ctx, isUpdate bool) (map[string]bool, error) {
	updatedFields := make(map[string]bool)

	if !isUpdate {
		u.ID = primitive.NewObjectID()
	}

	if name := c.FormValue("name"); name != "" {
		u.Name = name
		updatedFields["name"] = true
	} else if !isUpdate {
		return nil, fmt.Errorf("name is required for new users")
	}

	if password := c.FormValue("password"); password != "" {
		if err := u.SetPassword(password); err != nil {
			return nil, fmt.Errorf("failed to set password: %w", err)
		}
		updatedFields["password_hash"] = true
	} else if !isUpdate {
		return nil, fmt.Errorf("password is required for new users")
	}

	if isUpdate {
		file, err := c.FormFile("avatar")
		if err == nil {
			src, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open avatar file: %w", err)
			}
			defer src.Close()

			avatarData, err := io.ReadAll(src)
			if err != nil {
				return nil, fmt.Errorf("failed to read avatar file: %w", err)
			}

			u.AvatarData = avatarData
			u.AvatarType = file.Header.Get("Content-Type")
			updatedFields["avatar_data"] = true
			updatedFields["avatar_type"] = true
		} else if err != fiber.ErrUnprocessableEntity {
			return nil, fmt.Errorf("failed to process avatar file: %w", err)
		}
	} else {
		u.AvatarData = nil
		u.AvatarType = ""
	}

	now := time.Now()
	if !isUpdate {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
	updatedFields["updated_at"] = true

	return updatedFields, nil
}

type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	Token        string             `bson:"token"`
	CreatedAt    time.Time          `bson:"created_at"`
	ExpiresAt    time.Time          `bson:"expires_at"`
	LastActivity time.Time          `bson:"last_activity"`
}

type History struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Timestamp time.Time          `bson:"timestamp"`
}

type Post struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content    string             `json:"content" bson:"content"`
	AuthorID   primitive.ObjectID `json:"author_id" bson:"author_id"`
	LikesCount int                `json:"likes_count" bson:"likes_count"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

func (p *Post) Parse(c *fiber.Ctx, authorID primitive.ObjectID, isUpdate bool) error {
	if err := c.BodyParser(p); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if p.Content == "" {
		return fmt.Errorf("content is required")
	}

	if !isUpdate {
		p.ID = primitive.NewObjectID()
		p.CreatedAt = time.Now()
		p.LikesCount = 0

	}
	p.AuthorID = authorID

	p.UpdatedAt = time.Now()

	return nil
}

type NotificationType string

const (
	NotificationTypeLike NotificationType = "like"
	// Add other notification types here as needed
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      NotificationType   `bson:"type" json:"type"`
	PostID    primitive.ObjectID `bson:"post_id" json:"post_id"`
	LikedBy   primitive.ObjectID `bson:"liked_by" json:"liked_by"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

func NewNotificationLike(postID, likedBy primitive.ObjectID) *Notification {
	return &Notification{
		Type:      NotificationTypeLike,
		PostID:    postID,
		LikedBy:   likedBy,
		CreatedAt: time.Now(),
	}
}
