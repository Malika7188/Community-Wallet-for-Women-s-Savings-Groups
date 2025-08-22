package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stellar/go/keypair"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
)

func ContributeToRound(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		Round  int     `json:"round"`
		Amount float64 `json:"amount"`
		Secret string  `json:"secret"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// ✅ Validate the secret key belongs to the user
	kp, err := keypair.ParseFull(payload.Secret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid secret key format",
		})
	}

	if kp.Address() != user.Wallet {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Secret key does not match your wallet address",
		})
	}

	// Verify user is a member of the group
	var member models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, user.ID, "approved").First(&member).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not a member of this group"})
	}

	// Get group details
	var group models.Group
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}

	if group.Status != "active" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Group is not active"})
	}

	// Validate amount matches expected contribution
	if payload.Amount != group.ContributionAmount {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Amount must be %.2f XLM", group.ContributionAmount),
		})
	}

	// Check if user already contributed for this round
	var existingContribution models.RoundContribution
	if err := database.DB.Where("group_id = ? AND member_id = ? AND round = ?",
		groupID, member.ID, payload.Round).First(&existingContribution).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Already contributed for this round"})
	}

	// Perform direct XLM transfer from user to group wallet
	tx, err := services.SendXLM(payload.Secret, group.Wallet, fmt.Sprintf("%.7f", payload.Amount))
	if err != nil {
		fmt.Printf("❌ Failed to send XLM: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to transfer funds: %v", err),
		})
	}
	output := tx.Hash // Use the transaction hash as output
	fmt.Printf("✅ XLM transferred successfully. Transaction Hash: %s\n", output)

	// Note: If the smart contract still needs to track contributions,
	// you might need to call an unauthenticated 'record_contribution'
	// function on the smart contract here, passing the user's public key
	// and the amount. However, this would require modifying the smart contract
	// to have such a function and ensuring proper authorization for it.
	// For now, we are only performing the direct transfer.

	// Record the contribution
	contribution := models.RoundContribution{
		ID:        uuid.NewString(),
		GroupID:   groupID,
		MemberID:  member.ID,
		Round:     payload.Round,
		Amount:    payload.Amount,
		Status:    "confirmed",
		TxHash:    output,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&contribution).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Update or create round status
	if err := updateRoundStatus(groupID, payload.Round); err != nil {
		fmt.Printf("Warning: Failed to update round status: %v\n", err)
	}

	return c.JSON(fiber.Map{
		"message":      "Contribution successful",
		"contribution": contribution,
		"tx_hash":      output,
	})
}

func GetRoundStatus(c *fiber.Ctx) error {
	groupID := c.Params("id")
	round := c.QueryInt("round", 1)

	// Get round contributions
	var contributions []models.RoundContribution
	database.DB.Where("group_id = ? AND round = ?", groupID, round).
		Preload("Member").
		Preload("Member.User").
		Find(&contributions)

	// Get round status
	var roundStatus models.RoundStatus
	database.DB.Where("group_id = ? AND round = ?", groupID, round).First(&roundStatus)

	// Get all approved members for this group
	var allMembers []models.Member
	database.DB.Where("group_id = ? AND status = ?", groupID, "approved").
		Preload("User").
		Find(&allMembers)

	// Create contribution map for easy lookup
	contributionMap := make(map[string]models.RoundContribution)
	for _, contrib := range contributions {
		contributionMap[contrib.MemberID] = contrib
	}

	// Build response with member contribution status
	type MemberContributionStatus struct {
		Member       models.Member             `json:"member"`
		HasPaid      bool                      `json:"has_paid"`
		Contribution *models.RoundContribution `json:"contribution,omitempty"`
	}

	var memberStatuses []MemberContributionStatus
	for _, member := range allMembers {
		contrib, hasPaid := contributionMap[member.ID]
		status := MemberContributionStatus{
			Member:  member,
			HasPaid: hasPaid,
		}
		if hasPaid {
			status.Contribution = &contrib
		}
		memberStatuses = append(memberStatuses, status)
	}

	return c.JSON(fiber.Map{
		"round":         round,
		"round_status":  roundStatus,
		"member_status": memberStatuses,
		"total_members": len(allMembers),
		"paid_members":  len(contributions),
	})
}

func AuthorizeRoundPayout(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		Round int `json:"round"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Check if user is admin/creator
	var admin models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?",
		groupID, user.ID, []string{"creator", "admin"}).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only admins can authorize payouts"})
	}

	// Check if all members have contributed
	var group models.Group
	database.DB.First(&group, "id = ?", groupID)

	var totalMembers int64
	database.DB.Model(&models.Member{}).Where("group_id = ? AND status = ?", groupID, "approved").Count(&totalMembers)

	var contributionsCount int64
	database.DB.Model(&models.RoundContribution{}).Where("group_id = ? AND round = ? AND status = ?",
		groupID, payload.Round, "confirmed").Count(&contributionsCount)

	if contributionsCount < totalMembers {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Not all members have contributed. %d/%d paid", contributionsCount, totalMembers),
		})
	}

	// Update round status to authorized
	database.DB.Model(&models.RoundStatus{}).
		Where("group_id = ? AND round = ?", groupID, payload.Round).
		Updates(map[string]interface{}{
			"payout_authorized": true,
			"status":            "ready_for_payout",
		})

	// Get the recipient for this round from payout schedule
	var payoutSchedule models.PayoutSchedule
	if err := database.DB.Where("group_id = ? AND round = ?", groupID, payload.Round).
		Preload("Member").
		Preload("Member.User").
		First(&payoutSchedule).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payout schedule not found"})
	}

	// Notify all members about authorized payout
	var members []models.Member
	database.DB.Where("group_id = ? AND status = ?", groupID, "approved").Find(&members)

	for _, member := range members {
		services.CreateNotification(
			member.UserID,
			groupID,
			"round_payout_authorized",
			"Round Payout Authorized",
			fmt.Sprintf("Round %d payout has been authorized for %s", payload.Round, payoutSchedule.Member.User.Name),
		)
	}

	return c.JSON(fiber.Map{
		"message":   "Round payout authorized successfully",
		"round":     payload.Round,
		"recipient": payoutSchedule.Member.User.Name,
		"amount":    payoutSchedule.Amount,
	})
}

func updateRoundStatus(groupID string, round int) error {
	// Get total required amount and member count
	var group models.Group
	database.DB.First(&group, "id = ?", groupID)

	var totalMembers int64
	database.DB.Model(&models.Member{}).Where("group_id = ? AND status = ?", groupID, "approved").Count(&totalMembers)

	// Get current contributions for this round
	var contributionsCount int64
	var totalReceived float64
	database.DB.Model(&models.RoundContribution{}).
		Where("group_id = ? AND round = ? AND status = ?", groupID, round, "confirmed").
		Count(&contributionsCount)

	database.DB.Model(&models.RoundContribution{}).
		Where("group_id = ? AND round = ? AND status = ?", groupID, round, "confirmed").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalReceived)

	totalRequired := group.ContributionAmount * float64(totalMembers)

	// Update or create round status
	roundStatus := models.RoundStatus{
		GroupID:           groupID,
		Round:             round,
		TotalRequired:     totalRequired,
		TotalReceived:     totalReceived,
		ContributorsCount: int(contributionsCount),
		RequiredCount:     int(totalMembers),
		Status:            "collecting",
	}

	if contributionsCount >= totalMembers {
		roundStatus.Status = "ready_for_payout"
	}

	return database.DB.Where("group_id = ? AND round = ?", groupID, round).
		Assign(roundStatus).
		FirstOrCreate(&roundStatus).Error
}
