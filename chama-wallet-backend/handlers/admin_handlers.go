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

func NominateAdmin(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		NomineeID string `json:"nominee_id"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Check if nominator is a member
	var nominator models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, user.ID, "approved").First(&nominator).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not a group member"})
	}

	// Check if nominee is a member
	var nominee models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, payload.NomineeID, "approved").First(&nominee).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nominee is not a group member"})
	}

	// Check if already nominated
	var existing models.AdminNomination
	if database.DB.Where("group_id = ? AND nominee_id = ? AND status = ?",
		groupID, payload.NomineeID, "pending").First(&existing).Error == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already nominated"})
	}

	nomination := models.AdminNomination{
		ID:          uuid.NewString(),
		GroupID:     groupID,
		NominatorID: user.ID,
		NomineeID:   payload.NomineeID,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&nomination).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if nominee has 2 nominations, auto-approve as admin
	var nominationCount int64
	database.DB.Model(&models.AdminNomination{}).
		Where("group_id = ? AND nominee_id = ? AND status = ?", groupID, payload.NomineeID, "pending").
		Count(&nominationCount)

	if nominationCount >= 2 {
		// Update member role to admin
		database.DB.Model(&models.Member{}).
			Where("group_id = ? AND user_id = ?", groupID, payload.NomineeID).
			Update("role", "admin")

		// Update all nominations to approved
		database.DB.Model(&models.AdminNomination{}).
			Where("group_id = ? AND nominee_id = ?", groupID, payload.NomineeID).
			Update("status", "approved")

		// Send notification
		services.CreateNotification(
			payload.NomineeID,
			groupID,
			"admin_promotion",
			"Promoted to Admin",
			"You have been promoted to group admin",
		)
	}

	return c.JSON(fiber.Map{"message": "Nomination submitted successfully"})
}

func ApproveMember(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		MemberID string `json:"member_id"`
		Action   string `json:"action"` // approve or reject
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Check if user is admin/creator
	var admin models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?",
		groupID, user.ID, []string{"creator", "admin"}).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	status := "approved"
	if payload.Action == "reject" {
		status = "rejected"
	}

	if err := database.DB.Model(&models.Member{}).
		Where("id = ? AND group_id = ?", payload.MemberID, groupID).
		Update("status", status).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Send notification to member
	var member models.Member
	database.DB.Where("id = ?", payload.MemberID).Preload("User").First(&member)

	notificationType := "membership_approved"
	title := "Membership Approved"
	message := "Your group membership has been approved"

	if status == "rejected" {
		notificationType = "membership_rejected"
		title = "Membership Rejected"
		message = "Your group membership has been rejected"
	}

	services.CreateNotification(member.UserID, groupID, notificationType, title, message)

	// Check if group now meets minimum requirements and can be approved
	if status == "approved" {
		var group models.Group
		database.DB.First(&group, "id = ?", groupID)

		var approvedMemberCount int64
		database.DB.Model(&models.Member{}).Where("group_id = ? AND status = ?", groupID, "approved").Count(&approvedMemberCount)

		if approvedMemberCount >= int64(group.MinMembers) && !group.IsApproved {
			// Notify creator that group meets minimum requirements and can be approved
			services.CreateNotification(
				group.CreatorID,
				groupID,
				"group_ready",
				"Group Ready for Approval",
				fmt.Sprintf("Your group now has %d members and can be approved for activation", approvedMemberCount),
			)
		}
	}

	return c.JSON(fiber.Map{"message": "Member status updated successfully"})
}
