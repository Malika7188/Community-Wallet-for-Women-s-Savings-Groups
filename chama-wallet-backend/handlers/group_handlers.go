package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
	"chama-wallet-backend/config"
)

func ContributeToGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		From   string `json:"from"`
		Secret string `json:"secret"`
		Amount string `json:"amount"`
	}
	
	if err := c.BodyParser(&payload); err != nil {
		fmt.Printf("❌ Failed to parse request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if payload.From == "" || payload.Secret == "" || payload.Amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: from, secret, and amount are required",
		})
	}

	// Validate amount limits for mainnet
	if config.Config.IsMainnet {
		amount, err := strconv.ParseFloat(payload.Amount, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid amount format",
			})
		}
		
		minAmount := 0.0000001 // Minimum XLM amount
		if amount < minAmount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Amount below minimum of %f XLM", minAmount),
			})
		}
	}
	// Verify user is a member of the group
	var member models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, user.ID, "approved").First(&member).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not an approved member of this group",
		})
	}

	// Get group details
	group, err := services.GetGroupByID(groupID)
	if err != nil {
		fmt.Printf("❌ Group not found: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}

	// Validate group status
	if group.Status != "active" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Group is not active (current status: %s)", group.Status),
		})
	}

	// Validate contract ID
	if group.ContractID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Group contract ID is not set",
		})
	}

	// Validate user's wallet address matches the 'from' field
	if payload.From != user.Wallet {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "From address must match your wallet address",
		})
	}

	fmt.Printf("🔄 Processing contribution: %s XLM from %s to group %s (contract: %s) on %s\n", 
		payload.Amount, payload.From, group.Name, group.ContractID, config.Config.Network)

	// Make authenticated Soroban contract call
	output, err := services.ContributeWithAuth(group.ContractID, payload.From, payload.Amount, payload.Secret)
	if err != nil {
		fmt.Printf("❌ Soroban contribution failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Blockchain transaction failed: %v", err),
		})
	}

	// Record the contribution in the database
	contribution := models.Contribution{
		ID:        uuid.NewString(),
		GroupID:   groupID,
		UserID:    user.ID,
		Amount:    parseFloat64(payload.Amount),
		Status:    "confirmed",
		TxHash:    output,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&contribution).Error; err != nil {
		fmt.Printf("⚠️ Warning: Failed to record contribution in database: %v\n", err)
		// Don't fail the request since blockchain transaction succeeded
	}

	fmt.Printf("✅ Contribution successful on %s: %s\n", config.Config.Network, output)

	return c.JSON(fiber.Map{
		"message":      "Contribution successful",
		"group_id":     groupID,
		"group_name":   group.Name,
		"from":         payload.From,
		"to":           group.Wallet,
		"amount":       payload.Amount,
		"tx_hash":      output,
		"network":      config.Config.Network,
		"contribution": contribution,
	})
}

func GetUserGroups(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	groups, err := services.GetUserGroups(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(groups)
}

func GetNonGroupMembers(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Check if user is admin/creator of the group
	var admin models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?", 
		groupID, user.ID, []string{"creator", "admin"}).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	// Get users who are not members of this group
	var users []models.User
	err := database.DB.
		Where("id NOT IN (SELECT user_id FROM members WHERE group_id = ?)", groupID).
		Find(&users).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(users)
}

func InviteToGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Check if user is admin/creator
	var admin models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?", 
		groupID, user.ID, []string{"creator", "admin"}).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only admins can invite users"})
	}

	// Check if user exists
	var invitedUser models.User
	if err := database.DB.Where("email = ?", payload.Email).First(&invitedUser).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User with this email not found"})
	}

	// Check if user is already a member
	var existingMember models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ?", groupID, invitedUser.ID).First(&existingMember).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is already a member of this group"})
	}

	// Create invitation
	invitation := models.GroupInvitation{
		ID:        uuid.NewString(),
		GroupID:   groupID,
		InviterID: user.ID,
		Email:     payload.Email,
		UserID:    invitedUser.ID,
		Status:    "pending",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := database.DB.Create(&invitation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Create notification for invited user
	var group models.Group
	database.DB.First(&group, "id = ?", groupID)

	notification := models.Notification{
		ID:        uuid.NewString(),
		UserID:    invitedUser.ID,
		GroupID:   groupID,
		Type:      "group_invitation",
		Title:     "Group Invitation",
		Message:   fmt.Sprintf("You have been invited to join %s by %s", group.Name, user.Name),
		Status:    "unread",
		Read:      false,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&notification).Error; err != nil {
		fmt.Printf("⚠️ Failed to create notification: %v\n", err)
	} else {
		fmt.Printf("✅ Notification created for user %s\n", invitedUser.ID)
	}

	return c.JSON(fiber.Map{"message": "Invitation sent successfully"})
}

func ApproveGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Check if user is the creator
	var group models.Group
	if err := database.DB.Where("id = ? AND creator_id = ?", groupID, user.ID).First(&group).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only group creator can approve the group"})
	}

	// Check if group has minimum members
	var memberCount int64
	database.DB.Model(&models.Member{}).Where("group_id = ? AND status = ?", groupID, "approved").Count(&memberCount)
	if memberCount < int64(group.MinMembers) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Group needs at least %d members before approval (currently has %d)", group.MinMembers, memberCount),
		})
	}

	// Update group approval status
	if err := database.DB.Model(&group).Update("is_approved", true).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to approve group"})
	}

	// Notify all members that group is approved
	var members []models.Member
	database.DB.Where("group_id = ? AND status = ?", groupID, "approved").Find(&members)

	for _, member := range members {
		services.CreateNotification(
			member.UserID,
			groupID,
			"group_approved",
			"Group Approved",
			"Your group has been approved and is ready for activation",
		)
	}

	return c.JSON(fiber.Map{"message": "Group approved successfully"})
}

func ActivateGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		ContributionAmount float64  `json:"contribution_amount"`
		ContributionPeriod int      `json:"contribution_period"`
		PayoutOrder        []string `json:"payout_order"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Debug logging
	fmt.Printf("Received payout order: %+v\n", payload.PayoutOrder)

	// Check if user is admin/creator
	var admin models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?",
		groupID, user.ID, []string{"creator", "admin"}).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only admins can activate groups"})
	}

	// Check if group is approved
	var group models.Group
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}

	if !group.IsApproved {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Group must be approved before activation"})
	}

	// Validate payout order
	if len(payload.PayoutOrder) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payout order cannot be empty"})
	}

	// Convert payout order to JSON string
	payoutOrderJSON, err := json.Marshal(payload.PayoutOrder)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payout order"})
	}

	fmt.Printf("Payout order JSON: %s\n", string(payoutOrderJSON))

	// Calculate next contribution date
	nextContributionDate := time.Now().AddDate(0, 0, payload.ContributionPeriod)

	// Update group
	updates := map[string]interface{}{
		"status":                "active",
		"contribution_amount":   payload.ContributionAmount,
		"contribution_period":   payload.ContributionPeriod,
		"payout_order":         string(payoutOrderJSON),
		"current_round":        1,
		"next_contribution_date": nextContributionDate,
	}

	if err := database.DB.Model(&group).Where("id = ?", groupID).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Create payout schedule entries
	var members []models.Member
	database.DB.Where("group_id = ? AND status = ?", groupID, "approved").Find(&members)

	totalPayout := payload.ContributionAmount * float64(len(members))
	startDate := time.Now()

	for i, memberID := range payload.PayoutOrder {
		// Find the member
		var member models.Member
		for _, m := range members {
			if m.UserID == memberID {
				member = m
				break
			}
		}
		
		if member.ID == "" {
			continue // Skip if member not found
		}
		
		// Calculate due date for this round
		dueDate := startDate.AddDate(0, 0, i*payload.ContributionPeriod)
		
		schedule := models.PayoutSchedule{
			ID:        uuid.NewString(),
			GroupID:   groupID,
			MemberID:  member.ID,
			Round:     i + 1,
			Amount:    totalPayout,
			DueDate:   dueDate,
			Status:    "scheduled",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		database.DB.Create(&schedule)
	}

	return c.JSON(fiber.Map{"message": "Group activated successfully"})
}

// Helper function to parse float64 safely
func parseFloat64(s string) float64 {
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val
	}
	return 0.0
}

func JoinGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Check if group exists and is active
	var group models.Group
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}

	// Check if group is full
	var memberCount int64
	database.DB.Model(&models.Member{}).Where("group_id = ? AND status = ?", groupID, "approved").Count(&memberCount)
	if memberCount >= int64(group.MaxMembers) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Group is full"})
	}

	// Check if user is already a member
	var existingMember models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ?", groupID, user.ID).First(&existingMember).Error; err == nil {
		if existingMember.Status == "pending" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Join request already pending"})
		}
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Already a member of this group"})
	}

	// Create new member with pending status
	member := models.Member{
		ID:       uuid.NewString(),
		GroupID:  groupID,
		UserID:   user.ID,
		Wallet:   user.Wallet,
		Role:     "member",
		Status:   "pending",
		JoinedAt: time.Now(),
	}

	if err := database.DB.Create(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to join group"})
	}

	// Create notification for group admins
	var admins []models.Member
	database.DB.Where("group_id = ? AND role IN ? AND status = ?",
		groupID, []string{"creator", "admin"}, "approved").Find(&admins)

	for _, admin := range admins {
		notification := models.Notification{
			ID:        uuid.NewString(),
			UserID:    admin.UserID,
			GroupID:   groupID,
			Type:      "join_request",
			Title:     "New Join Request",
			Message:   fmt.Sprintf("%s wants to join %s", user.Name, group.Name),
			Data:      fmt.Sprintf(`{"group_id":"%s","member_id":"%s"}`, groupID, member.ID),
			Status:    "unread",
			CreatedAt: time.Now(),
		}
		database.DB.Create(&notification)
	}

	return c.JSON(fiber.Map{
		"message": "Join request sent successfully",
		"member":  member,
	})
}

func GetPayoutSchedule(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Check if user is member of the group
	var member models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, user.ID, "approved").First(&member).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not a group member"})
	}

	var schedules []models.PayoutSchedule
	err := database.DB.Where("group_id = ?", groupID).
		Preload("Member.User").
		Order("round ASC").
		Find(&schedules).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(schedules)
}
