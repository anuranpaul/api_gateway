package handlers

import (
	"example/API_Gateway/internal/models"
	"example/API_Gateway/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
    userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
    return &UserHandler{userRepo: userRepo}
}

// CreateUser handles user creation
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req models.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Check if requester is admin for creating admin users
    if req.Role == "admin" {
        userRole := c.Locals("role").(string)
        if userRole != "admin" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Only admins can create admin users",
            })
        }
    }

    user, err := h.userRepo.CreateUser(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUser handles retrieving a single user
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid user ID",
        })
    }

    // Check if user is requesting their own data or is an admin
    userRole := c.Locals("role").(string)
    userID := c.Locals("user_id").(int)
    if userRole != "admin" && userID != id {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Access denied",
        })
    }

    user, err := h.userRepo.GetUserByID(c.Context(), id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User not found",
        })
    }

    return c.JSON(user)
}

// GetAllUsers handles retrieving all users (admin only)
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
    // Verify admin role
    if c.Locals("role").(string) != "admin" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Admin access required",
        })
    }

    users, err := h.userRepo.GetAllUsers(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(users)
}

// UpdateUser handles user updates
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid user ID",
        })
    }

    var req models.UpdateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Check if user is updating their own data or is an admin
    userRole := c.Locals("role").(string)
    userID := c.Locals("user_id").(int)
    if userRole != "admin" && userID != id {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Access denied",
        })
    }

    // Only admins can change roles
    if req.Role != "" && userRole != "admin" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Only admins can change roles",
        })
    }

    user, err := h.userRepo.UpdateUser(c.Context(), id, &req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(user)
}

// DeleteUser handles user deletion
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid user ID",
        })
    }

    // Check if user is deleting their own account or is an admin
    userRole := c.Locals("role").(string)
    userID := c.Locals("user_id").(int)
    if userRole != "admin" && userID != id {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Access denied",
        })
    }

    if err := h.userRepo.DeleteUser(c.Context(), id); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.SendStatus(fiber.StatusNoContent)
} 