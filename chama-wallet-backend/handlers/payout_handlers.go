package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
)

func CreatePayoutRequest(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		RecipientID string  `json:"recipient_id"`
		Amount      float64 `json:"amount"`
		Round       int     `json:"round"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Check if user is admin/creator
	var admin models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?",
		groupID, user.ID, []string{"creator", "admin"}).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only admins can create payout requests"})
	}

	// Verify recipient is group member
	var recipient models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, payload.RecipientID, "approved").First(&recipient).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Recipient is not a group member"})
	}

	payoutRequest := models.PayoutRequest{
		ID:          uuid.NewString(),
		GroupID:     groupID,
		RecipientID: payload.RecipientID,
		Amount:      payload.Amount,
		Round:       payload.Round,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&payoutRequest).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Notify all admins about the payout request
	var admins []models.Member
	database.DB.Where("group_id = ? AND role IN ? AND status = ?",
		groupID, []string{"creator", "admin"}, "approved").Find(&admins)

	for _, admin := range admins {
		if admin.UserID != user.ID { // Don't notify the creator
			services.CreateNotification(
				admin.UserID,
				groupID,
				"payout_request",
				"Payout Request Created",
				"A new payout request requires your approval",
			)
		}
	}

	return c.JSON(fiber.Map{
		"message": "Payout request created successfully",
		"request": payoutRequest,
	})
}

func ApprovePayoutRequest(c *fiber.Ctx) error {
	payoutID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		Approved bool `json:"approved"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Get payout request
	var payoutRequest models.PayoutRequest
	if err := database.DB.Where("id = ?", payoutID).Preload("Group").First(&payoutRequest).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payout request not found"})
	}

	// Check if user is admin/creator of the group
	var admin models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?",
		payoutRequest.GroupID, user.ID, []string{"creator", "admin"}).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	// Check if already approved by this admin
	var existingApproval models.PayoutApproval
	if database.DB.Where("payout_request_id = ? AND admin_id = ?",
		payoutID, user.ID).First(&existingApproval).Error == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Already voted on this request"})
	}

	// Create approval record
	approval := models.PayoutApproval{
		ID:              uuid.NewString(),
		PayoutRequestID: payoutID,
		AdminID:         user.ID,
		Approved:        payload.Approved,
		CreatedAt:       time.Now(),
	}

	if err := database.DB.Create(&approval).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if we have enough approvals (2 admins)
	var approvalCount int64
	database.DB.Model(&models.PayoutApproval{}).
		Where("payout_request_id = ? AND approved = ?", payoutID, true).
		Count(&approvalCount)

	var rejectionCount int64
	database.DB.Model(&models.PayoutApproval{}).
		Where("payout_request_id = ? AND approved = ?", payoutID, false).
		Count(&rejectionCount)

	if approvalCount >= 2 {
		// Approve the payout
		database.DB.Model(&models.PayoutRequest{}).
			Where("id = ?", payoutID).
			Update("status", "approved")

		// Notify all members about approved payout
		var members []models.Member
		database.DB.Where("group_id = ? AND status = ?",
			payoutRequest.GroupID, "approved").Find(&members)

		for _, member := range members {
			services.CreateNotification(
				member.UserID,
				payoutRequest.GroupID,
				"payout_approved",
				"Payout Approved",
				"A payout has been approved and will be processed",
			)
		}

		return c.JSON(fiber.Map{"message": "Payout request approved and will be processed"})
	} else if rejectionCount >= 1 {
		// Reject the payout
		database.DB.Model(&models.PayoutRequest{}).
			Where("id = ?", payoutID).
			Update("status", "rejected")

		return c.JSON(fiber.Map{"message": "Payout request rejected"})
	}

	return c.JSON(fiber.Map{"message": "Approval recorded, waiting for more approvals"})
}

func GetPayoutRequests(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Check if user is member of the group
	var member models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, user.ID, "approved").First(&member).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not a group member"})
	}

	var payoutRequests []models.PayoutRequest
	err := database.DB.Where("group_id = ?", groupID).
		Preload("Recipient").
		Preload("Approvals.Admin").
		Order("created_at DESC").
		Find(&payoutRequests).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(payoutRequests)
}
