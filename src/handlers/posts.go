package handlers

import (
	"glamapp/src/database"
	"glamapp/src/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePost handles the creation of a new post
func CreatePost(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		post := new(models.Post)
		if err := c.BodyParser(post); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "Invalid request body",
			})
		}

		// Convert the author ID from string to ObjectID
		authorID, err := primitive.ObjectIDFromHex(c.FormValue("author_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "Invalid author ID",
			})
		}
		post.AuthorID = authorID

		id, err := db.CreatePost(post)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": "Failed to create post",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id": id,
		})
	}
}

// GetPost handles fetching a single post by its ID
func GetPost(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		post, err := db.GetPost(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"detail": "Post not found",
			})
		}

		return c.JSON(fiber.Map{
			"data": post,
		})
	}
}

// UpdatePost handles updating an existing post
func UpdatePost(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		post := new(models.Post)
		if err := c.BodyParser(post); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "Invalid request body",
			})
		}

		if err := db.UpdatePost(id, post); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": "Failed to update post",
			})
		}

		return c.JSON(fiber.Map{
			"detail": "Post updated successfully",
		})
	}
}

// DeletePost handles deleting a post
func DeletePost(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := db.DeletePost(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": "Failed to delete post",
			})
		}

		return c.JSON(fiber.Map{
			"detail": "Post deleted successfully",
		})
	}
}

// GetPostsByProfile handles fetching all posts by a specific profile
func GetPostsByProfile(db *database.MongoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		profileID := c.Params("profileId")
		posts, err := db.GetPostsByProfile(profileID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": "Failed to fetch posts",
			})
		}

		return c.JSON(fiber.Map{
			"data":  posts,
			"count": len(posts),
		})
	}
}
