package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
)

func GetNotifications(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	fmt.Printf("🔍 Getting notifications for user: %s\n", user.ID)

	notifications, err := services.GetUserNotifications(user.ID)
	if err != nil {
		fmt.Printf("❌ Error getting notifications: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("✅ Found %d notifications for user %s\n", len(notifications), user.ID)
	return c.JSON(notifications)
}

func MarkNotificationRead(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Verify notification belongs to user
	var notification models.Notification
	if err := database.DB.Where("id = ? AND user_id = ?", notificationID, user.ID).First(&notification).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Notification not found"})
	}

	if err := services.MarkNotificationAsRead(notificationID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Notification marked as read"})
}

func GetUserInvitations(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	fmt.Printf("🔍 Getting invitations for user: %s (email: %s)\n", user.ID, user.Email)

	var invitations []models.GroupInvitation
	err := database.DB.Where("email = ? AND status = ?", user.Email, "pending").
		Preload("Group").
		Preload("Inviter").
		Find(&invitations).Error

	if err != nil {
		fmt.Printf("❌ Error getting invitations: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("✅ Found %d invitations for user %s\n", len(invitations), user.Email)
	return c.JSON(invitations)
}

func AcceptInvitation(c *fiber.Ctx) error {
	invitationID := c.Params("id")
	user := c.Locals("user").(models.User)

	var invitation models.GroupInvitation
	if err := database.DB.Where("id = ? AND email = ?", invitationID, user.Email).First(&invitation).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invitation not found"})
	}

	if invitation.Status != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invitation already processed"})
	}

	if time.Now().After(invitation.ExpiresAt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invitation has expired"})
	}

	// Update invitation status
	database.DB.Model(&invitation).Updates(map[string]interface{}{
		"status":  "accepted",
		"user_id": user.ID,
	})

	// Add user as member with approved status (since they were invited)
	member := models.Member{
		ID:       uuid.NewString(),
		GroupID:  invitation.GroupID,
		UserID:   user.ID,
		Wallet:   user.Wallet,
		Role:     "member",
		Status:   "approved", // Auto-approve invited users
		JoinedAt: time.Now(),
	}

	if err := database.DB.Create(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Notify group admins about new member (not request)
	var admins []models.Member
	database.DB.Where("group_id = ? AND role IN ? AND status = ?",
		invitation.GroupID, []string{"creator", "admin"}, "approved").Find(&admins)

	for _, admin := range admins {
		services.CreateNotification(
			admin.UserID,
			invitation.GroupID,
			"new_member_joined",
			"New Member Joined",
			user.Name+" has joined the group",
		)
	}

	// Notify the user that they successfully joined
	services.CreateNotification(
		user.ID,
		invitation.GroupID,
		"membership_approved",
		"Welcome to the Group",
		"You have successfully joined the group",
	)

	return c.JSON(fiber.Map{"message": "Invitation accepted successfully"})
}

func RejectInvitation(c *fiber.Ctx) error {
	invitationID := c.Params("id")
	user := c.Locals("user").(models.User)

	var invitation models.GroupInvitation
	if err := database.DB.Where("id = ? AND email = ?", invitationID, user.Email).First(&invitation).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invitation not found"})
	}

	if invitation.Status != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invitation already processed"})
	}

	// Update invitation status
	database.DB.Model(&invitation).Update("status", "rejected")

	return c.JSON(fiber.Map{"message": "Invitation rejected"})
}
