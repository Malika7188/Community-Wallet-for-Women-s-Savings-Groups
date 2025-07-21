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

func ContributeToGroup(c *fiber.Ctx) error {
	// ContractID := "CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4"

	groupID := c.Params("id")

	var payload struct {
		From   string `json:"from"`
		Secret string `json:"secret"`
		Amount string `json:"amount"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	group, err := services.GetGroupByID(groupID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "group not found"})
	}

	// 🔁 Soroban contract call instead of native XLM payment
	args := []string{payload.From, payload.Amount}
	output, err := services.CallSorobanFunction(group.ContractID, "contribute", args)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Contribution successful via Soroban",
		"group":   groupID,
		"from":    payload.From,
		"to":      group.Wallet,
		"amount":  payload.Amount,
		"tx":      output,
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
	
	services.CreateNotification(
		invitedUser.ID,
		groupID,
		"group_invitation",
		"Group Invitation",
		fmt.Sprintf("You have been invited to join %s by %s", group.Name, user.Name),
	)

	return c.JSON(fiber.Map{"message": "Invitation sent successfully"})
}

func GetNonGroupMembers(c *fiber.Ctx) error {
	groupID := c.Params("id")

	users, err := services.GetNonGroupMembers(groupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(users)
}

func ApproveGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Check if user is the creator
	var group models.Group
	if err := database.DB.Where("id = ? AND creator_id = ?", groupID, user.ID).First(&group).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only group creator can approve the group"})
	}

	// Check if group is full
	var memberCount int64
	database.DB.Model(&models.Member{}).Where("group_id = ? AND status = ?", groupID, "approved").Count(&memberCount)
	if memberCount < int64(group.MaxMembers) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Group must be full before approval"})
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
		ContributionAmount float64 `json:"contribution_amount"`
		ContributionPeriod int     `json:"contribution_period"`
		PayoutOrder        []string `json:"payout_order"`
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

	// Check if group is approved
	var group models.Group
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}

	if !group.IsApproved {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Group must be approved before activation"})
	}

	// Calculate next contribution date
	nextContributionDate := time.Now().AddDate(0, 0, payload.ContributionPeriod)

	// Convert payout order to JSON string
	payoutOrderJSON, err := json.Marshal(payload.PayoutOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process payout order"})
	}

	// Update group with activation settings
	updates := map[string]interface{}{
		"status": "active",
		"contribution_amount": payload.ContributionAmount,
		"contribution_period": payload.ContributionPeriod,
		"payout_order": string(payoutOrderJSON),
		"next_contribution_date": nextContributionDate,
	}

	if err := database.DB.Model(&models.Group{}).Where("id = ?", groupID).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Notify all members about activation
	var members []models.Member
	database.DB.Where("group_id = ? AND status = ?", groupID, "approved").Find(&members)

	for _, member := range members {
		services.CreateNotification(
			member.UserID,
			groupID,
			"group_activated",
			"Group Activated",
			fmt.Sprintf("Your group is now active! Contribution amount: %.2f XLM every %d days", payload.ContributionAmount, payload.ContributionPeriod),
		)
	}

	return c.JSON(fiber.Map{"message": "Group activated successfully"})
}

func JoinGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Use user's wallet address
	walletAddress := user.Wallet

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
		Wallet:   walletAddress,
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
