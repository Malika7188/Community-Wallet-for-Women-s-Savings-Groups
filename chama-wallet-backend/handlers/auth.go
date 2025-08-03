package handlers

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
)

// Register handles user registration
func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, email, and password are required",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters long",
		})
	}

	// Register user
	authResponse, err := services.RegisterUser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user": fiber.Map{
			"id":         authResponse.User.ID,
			"name":       authResponse.User.Name,
			"email":      authResponse.User.Email,
			"wallet":     authResponse.User.Wallet,
			"secret_key": authResponse.User.SecretKey, // Include secret key in response
		},
		"token": authResponse.Token,
	})
}

// Login handles user authentication
func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Login user
	authResponse, err := services.LoginUser(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(authResponse)
}

// GetProfile returns the current user's profile including secret key
func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	user, err := services.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"wallet":     user.Wallet,
			"secret_key": user.SecretKey, // Add this to show user their secret key
			"created_at": user.CreatedAt,
		},
	})
}

// UpdateProfile updates the current user's profile
func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get current user
	user, err := services.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update name if provided
	if req.Name != "" {
		user.Name = req.Name
	}

	// Save changes
	if err := services.UpdateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// Logout handles user logout (client-side token removal)
func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
