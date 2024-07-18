package handlers

import (
	"glamapp/src/database"
	"glamapp/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type PostHandler struct {
	DB     *database.MongoDB
	Logger zerolog.Logger
}

func NewPostHandler(db *database.MongoDB, logger zerolog.Logger) *PostHandler {
	return &PostHandler{
		DB:     db,
		Logger: logger,
	}
}

// CreatePost handles the creation of a new post
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	var post models.Post
	if err := post.Parse(c, false); err != nil {
		h.Logger.Error().Err(err).Msg("Failed to parse post data")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id, err := h.DB.CreatePost(&post)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to create post")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create post")
	}

	h.Logger.Info().Str("id", id).Msg("Post created successfully")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post
	if err := post.Parse(c, true); err != nil {
		h.Logger.Error().Err(err).Str("id", id).Msg("Failed to parse post update data")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.DB.UpdatePost(id, &post); err != nil {
		h.Logger.Error().Err(err).Str("id", id).Msg("Failed to update post")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update post")
	}

	h.Logger.Info().Str("id", id).Msg("Post updated successfully")
	return c.JSON(fiber.Map{
		"detail": "Post updated successfully",
	})
}

func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	id := c.Params("id")
	post, err := h.DB.GetPost(id)
	if err != nil {
		h.Logger.Error().Err(err).Str("id", id).Msg("Failed to fetch post")
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}

	return c.JSON(fiber.Map{
		"data": post,
	})
}

func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.DB.DeletePost(id); err != nil {
		h.Logger.Error().Err(err).Str("id", id).Msg("Failed to delete post")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete post")
	}

	h.Logger.Info().Str("id", id).Msg("Post deleted successfully")
	return c.JSON(fiber.Map{
		"detail": "Post deleted successfully",
	})
}

func (h *PostHandler) GetPostsByProfile(c *fiber.Ctx) error {
	profileID := c.Params("profileId")
	posts, err := h.DB.GetPostsByProfile(profileID)
	if err != nil {
		h.Logger.Error().Err(err).Str("profileId", profileID).Msg("Failed to fetch posts")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch posts")
	}

	return c.JSON(fiber.Map{
		"data":  posts,
		"count": len(posts),
	})
}
