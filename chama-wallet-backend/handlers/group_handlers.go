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

func InviteToGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	if err := services.InviteUserToGroup(groupID, user.ID, payload.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

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

func ActivateGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		ContributionAmount float64 `json:"contribution_amount"`
		ContributionPeriod int     `json:"contribution_period"`
		PayoutOrder        string  `json:"payout_order"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	settings := models.GroupSettings{
		ContributionAmount: payload.ContributionAmount,
		ContributionPeriod: payload.ContributionPeriod,
		PayoutOrder:        payload.PayoutOrder,
	}

	if err := services.ApproveGroupActivation(groupID, user.ID, settings); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Group activated successfully"})
}

func JoinGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		Wallet string `json:"wallet"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Use user's wallet if not provided
	walletAddress := payload.Wallet
	if walletAddress == "" {
		walletAddress = user.Wallet
	}

	// Check if group exists and is active
	var group models.Group
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}

	// Check if user is already a member
	var existingMember models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ?", groupID, user.ID).First(&existingMember).Error; err == nil {
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
