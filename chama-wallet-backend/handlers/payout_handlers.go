package handlers

import (
	"fmt"
	"strconv"
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
		fmt.Printf("❌ Failed to parse payout request: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Validate required fields
	if payload.RecipientID == "" || payload.Amount <= 0 || payload.Round <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid required fields",
		})
	}

	// Get group details
	var group models.Group
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}

	// Validate group is active
	if group.Status != "active" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Group must be active to create payout requests",
		})
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

	// Check if payout request already exists for this round
	var existingRequest models.PayoutRequest
	if err := database.DB.Where("group_id = ? AND round = ? AND status IN ?",
		groupID, payload.Round, []string{"pending", "approved"}).First(&existingRequest).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Payout request already exists for this round",
		})
	}

	// Validate payout amount against group balance
	groupBalance, err := services.CheckBalance(group.Wallet)
	if err != nil {
		fmt.Printf("⚠️ Warning: Could not check group balance: %v\n", err)
	} else {
		if balance, parseErr := strconv.ParseFloat(groupBalance, 64); parseErr == nil {
			if payload.Amount > balance {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Insufficient group balance. Available: %.2f XLM, Requested: %.2f XLM", balance, payload.Amount),
				})
			}
		}
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
		fmt.Printf("❌ Failed to create payout request: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("✅ Payout request created: %s\n", payoutRequest.ID)

	// Notify all admins about the payout request (excluding the creator)
	var admins []models.Member
	database.DB.Where("group_id = ? AND role IN ? AND status = ?",
		groupID, []string{"creator", "admin"}, "approved").Find(&admins)

	for _, admin := range admins {
		if admin.UserID != user.ID {
			services.CreateNotification(
				admin.UserID,
				groupID,
				"payout_request",
				"Payout Request Created",
				fmt.Sprintf("New payout request for %.2f XLM to %s requires approval", payload.Amount, recipient.User.Name),
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
		fmt.Printf("❌ Failed to parse approval request: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Get payout request
	var payoutRequest models.PayoutRequest
	if err := database.DB.Where("id = ?", payoutID).Preload("Group").First(&payoutRequest).Error; err != nil {
		fmt.Printf("❌ Payout request not found: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payout request not found"})
	}

	// Check if payout is still pending
	if payoutRequest.Status != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Payout request is already %s", payoutRequest.Status),
		})
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
		fmt.Printf("❌ Failed to create approval record: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("✅ Approval recorded: %s by %s\n", payoutID, user.Name)

	// Count approvals and rejections
	var approvalCount int64
	database.DB.Model(&models.PayoutApproval{}).
		Where("payout_request_id = ? AND approved = ?", payoutID, true).
		Count(&approvalCount)

	var rejectionCount int64
	database.DB.Model(&models.PayoutApproval{}).
		Where("payout_request_id = ? AND approved = ?", payoutID, false).
		Count(&rejectionCount)

	fmt.Printf("📊 Payout %s - Approvals: %d, Rejections: %d\n", payoutID, approvalCount, rejectionCount)

	// Process payout if we have enough approvals (1 admins) or any rejection
	if approvalCount >= 1 {
		fmt.Printf("✅ Payout approved with %d approvals, processing...\n", approvalCount)
		
		// Update payout status to approved
		database.DB.Model(&models.PayoutRequest{}).
			Where("id = ?", payoutID).
			Update("status", "approved")

		// Execute the actual payout using Soroban contract
		if err := executePayout(payoutRequest); err != nil {
			fmt.Printf("❌ Payout execution failed: %v\n", err)
			// Update status to failed
			database.DB.Model(&models.PayoutRequest{}).
				Where("id = ?", payoutID).
				Update("status", "failed")
			
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Payout execution failed: %v", err),
			})
		}

		// Update status to completed
		database.DB.Model(&models.PayoutRequest{}).
			Where("id = ?", payoutID).
			Update("status", "completed")

		// Fetch group details again to get the latest CurrentRound and ContributionPeriod
		var currentGroup models.Group
		if err := database.DB.First(&currentGroup, "id = ?", payoutRequest.GroupID).Error; err != nil {
			fmt.Printf("❌ Failed to fetch group details for round update: %v\n", err)
		} else {
			// Increment current round and update next contribution date
			database.DB.Model(&models.Group{}).
				Where("id = ?", payoutRequest.GroupID).
				Update("current_round", currentGroup.CurrentRound + 1).
				Update("next_contribution_date", time.Now().AddDate(0, 0, currentGroup.ContributionPeriod))
		}

		// Notify all members about successful payout
		var members []models.Member
		database.DB.Where("group_id = ? AND status = ?",
			payoutRequest.GroupID,
			"approved").Find(&members)

		for _, member := range members {
			services.CreateNotification(
				member.UserID,
				payoutRequest.GroupID,
				"payout_approved",
				"Payout Approved",
				fmt.Sprintf("Payout of %.2f XLM has been approved and processed", payoutRequest.Amount),
			)
		}

		return c.JSON(fiber.Map{
			"message": "Payout request approved and executed successfully",
			"status":  "completed",
		})
	} else if rejectionCount >= 1 {
		fmt.Printf("❌ Payout rejected with %d rejections\n", rejectionCount)
		
		database.DB.Model(&models.PayoutRequest{}).
			Where("id = ?", payoutID).
			Update("status", "rejected")

		return c.JSON(fiber.Map{
			"message": "Payout request rejected",
			"status":  "rejected",
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Approval recorded, waiting for more approvals (%d/1)", approvalCount),
		"status":  "pending",
	})
}

// executePayout performs the actual blockchain payout transaction
func executePayout(payoutRequest models.PayoutRequest) error {
	fmt.Printf("🔄 Executing payout: %.2f XLM to recipient %s\n", payoutRequest.Amount, payoutRequest.RecipientID)
	
	// Get group details
	var group models.Group
	if err := database.DB.First(&group, "id = ?", payoutRequest.GroupID).Error; err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	// Get recipient details
	var recipient models.User
	if err := database.DB.First(&recipient, "id = ?", payoutRequest.RecipientID).Error; err != nil {
		return fmt.Errorf("failed to get recipient: %w", err)
	}

	// Validate group has secret key for transactions
	if group.SecretKey == "" {
		return fmt.Errorf("group secret key not available")
	}

	// Send actual XLM from group wallet to recipient
	_, err := services.SendXLM(group.SecretKey, recipient.Wallet, fmt.Sprintf("%.7f", payoutRequest.Amount))
	if err != nil {
		fmt.Printf("⚠️ Warning: XLM transfer failed but contract withdrawal succeeded: %v\n", err)
		return fmt.Errorf("soroban withdrawal failed: %w", err)
	}

	fmt.Printf("✅ Payout executed successfully")
	return nil
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
