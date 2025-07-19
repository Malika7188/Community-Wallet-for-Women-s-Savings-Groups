package handlers

import (
	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetNotifications(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	notifications, err := services.GetUserNotifications(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

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

	var invitations []models.GroupInvitation
	err := database.DB.Where("email = ? AND status = ?", user.Email, "pending").
		Preload("Group").
		Preload("Inviter").
		Find(&invitations).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

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

	// Add user as member
	member := models.Member{
		ID:       uuid.NewString(),
		GroupID:  invitation.GroupID,
		UserID:   user.ID,
		Wallet:   user.Wallet,
		Role:     "member",
		Status:   "pending", // Requires admin approval
		JoinedAt: time.Now(),
	}

	if err := database.DB.Create(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Notify group admins
	var admins []models.Member
	database.DB.Where("group_id = ? AND role IN ? AND status = ?",
		invitation.GroupID, []string{"creator", "admin"}, "approved").Find(&admins)

	for _, admin := range admins {
		services.CreateNotification(
			admin.UserID,
			invitation.GroupID,
			"new_member_request",
			"New Member Request",
			user.Name+" has accepted an invitation and requests to join the group",
		)
	}

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
