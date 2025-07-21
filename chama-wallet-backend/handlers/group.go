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

type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateGroup(c *fiber.Ctx) error {
	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Get authenticated user
	user := c.Locals("user").(models.User)

	// Generate wallet for the group
	wallet, err := services.GenerateStellarWallet()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate wallet"})
	}

	// ✅ Step 1: Deploy contract using CLI (or pre-deployed if needed)
	contractID, err := services.DeployChamaContract()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to deploy contract"})
	}

	// ✅ Step 2: Save group in DB with the new contractID
	group := models.Group{
		ID:          uuid.NewString(),
		Name:        payload.Name,
		Description: payload.Description,
		Wallet:      wallet.PublicKey,
		CreatorID:   user.ID,
		ContractID:  contractID,
		Status:      "pending",
	}

	if err := database.DB.Create(&group).Error; err != nil {
		fmt.Printf("❌ Failed to create group: %v\n", err) // Add debug log
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save group"})
	}

	fmt.Printf("✅ Group created successfully: %+v\n", group) // Add debug log

	// Add the creator as the first member with creator role and approved status
	member := models.Member{
		ID:       uuid.NewString(),
		GroupID:  group.ID,
		UserID:   user.ID,
		Wallet:   user.Wallet,
		Role:     "creator",
		Status:   "approved",
		JoinedAt: time.Now(),
	}

	err = database.DB.Create(&member).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add creator as member"})
	}

	return c.JSON(fiber.Map{
		"message":     "Group created",
		"group":       group,
		"contract_id": contractID,
	})
}

func AddMember(c *fiber.Ctx) error {
	groupID := c.Params("id")

	var body struct {
		Wallet string `json:"wallet"`
		UserID string `json:"user_id"` // Add this field
	}

	if err := c.BodyParser(&body); err != nil || body.Wallet == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request. Wallet is required.",
		})
	}

	group, err := services.AddMemberToGroup(groupID, body.UserID, body.Wallet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(group)
}
func DepositToGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")

	var body struct {
		FromWallet string `json:"from_wallet"`
		Secret     string `json:"secret"` // sender's secret key
		Amount     string `json:"amount"` // XLM to deposit
	}

	if err := c.BodyParser(&body); err != nil || body.FromWallet == "" || body.Secret == "" || body.Amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields.",
		})
	}

	group, err := services.GetGroupByID(groupID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Group not found.",
		})
	}

	// Send XLM to group wallet
	tx, err := services.SendXLM(body.Secret, group.Wallet, body.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":   "Deposit successful",
		"tx_id":     tx.ID,
		"from":      body.FromWallet,
		"to":        group.Wallet,
		"amount":    body.Amount,
		"timestamp": tx.LedgerCloseTime,
	})
}

// GetGroupBalance returns the XLM balance of the group's wallet
func GetGroupBalance(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if the group exists
	group, err := services.GetGroupByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Group not found",
		})
	}

	// Fetch the balance from Stellar
	balance, err := services.GetBalance(group.Wallet, group.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return group wallet and balance
	return c.JSON(fiber.Map{
		"group_id": id,
		"wallet":   group.Wallet,
		"balance":  balance,
	})
}
func GetAllGroups(c *fiber.Ctx) error {
	// Get authenticated user (optional)
	user := c.Locals("user")

	var groups []models.Group
	var err error

	if user != nil {
		// If user is authenticated, only show groups they are part of
		userModel := user.(models.User)
		groups, err = services.GetUserGroups(userModel.ID)
	} else {
		// If not authenticated, show no groups
		groups = []models.Group{}
	}

	if err != nil {
		fmt.Printf("❌ Error fetching groups: %v\n", err) // Add debug log
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch groups",
		})
	}

	fmt.Printf("✅ Found %d groups\n", len(groups)) // Add debug log
	return c.JSON(groups)
}

func GetGroupDetails(c *fiber.Ctx) error {
	groupID := c.Params("id")
	group, err := services.GetGroupWithMembers(groupID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Group not found",
		})
	}
	return c.JSON(group)
}
