package handlers

import (
	"glamapp/src/database"
	"glamapp/src/models"
	"strings"

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

func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var post models.Post
	if err := post.Parse(c, user.ID, false); err != nil {
		h.Logger.Error().Err(err).Msg("Failed to parse post data")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id, err := h.DB.CreatePost(&post, user.ID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to create post")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create post")
	}

	h.Logger.Info().Str("id", id).Str("user_id", user.ID.Hex()).Msg("Post created successfully")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func (h *PostHandler) LikePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	user := c.Locals("user").(*models.User)

	err := h.DB.LikePost(postID, user.ID)
	if err != nil {
		switch {
		case err.Error() == "post already liked or not found":
			h.Logger.Warn().Err(err).Str("post_id", postID).Str("user_id", user.ID.Hex()).Msg("Post already liked or not found")
			return fiber.NewError(fiber.StatusBadRequest, "Post already liked or not found")
		case strings.Contains(err.Error(), "failed to update user"):
			h.Logger.Error().Err(err).Str("post_id", postID).Str("user_id", user.ID.Hex()).Msg("Failed to update user's liked posts")
			// The post was liked, but we couldn't update the user's liked_posts list
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"detail": "Post liked successfully, but there was an issue updating your profile. Please refresh.",
			})
		default:
			h.Logger.Error().Err(err).Str("post_id", postID).Str("user_id", user.ID.Hex()).Msg("Failed to like post")
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to like post")
		}
	}

	h.Logger.Info().Str("post_id", postID).Str("user_id", user.ID.Hex()).Msg("Post liked successfully")
	return c.JSON(fiber.Map{
		"detail": "Post liked successfully",
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

func (h *PostHandler) GetPosts(c *fiber.Ctx) error {
	posts, err := h.DB.GetPosts()
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to fetch posts")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch posts")
	}

	return c.JSON(fiber.Map{
		"data":  posts,
		"count": len(posts),
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
