package handlers

import (
	"chama-wallet-backend/utils"
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
	wallet, err := utils.GenerateStellarWallet()
	if err != nil {
		fmt.Printf("❌ Failed to generate group wallet: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate wallet"})
	}

	fmt.Printf("✅ Generated group wallet: %s\n", wallet.PublicKey)

	// Remove the funding section - comment out or delete these lines:
	// err = services.FundTestAccount(wallet.PublicKey)
	// if err != nil {
	//     fmt.Printf("⚠️ Warning: Failed to fund group wallet: %v\n", err)
	//     // Don't fail the group creation, just log the warning
	// } else {
	//     fmt.Printf("✅ Group wallet funded successfully\n")
	// }

	// Deploy contract (non-blocking - don't fail group creation if contract deployment fails)
	contractID, err := services.DeployChamaContract()
	if err != nil {
		fmt.Printf("⚠️ Warning: Failed to deploy contract: %v\n", err)
		contractID = "" // Set empty contract ID, can be updated later
		// Don't fail the group creation - contract can be deployed later
	} else {
		fmt.Printf("✅ Contract deployed successfully: %s\n", contractID)
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
		SecretKey:   wallet.SecretKey, // Add this field to store secret
	}

	if err := database.DB.Create(&group).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

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
		fmt.Printf("⚠️ Warning: Failed to add creator as member: %v\n", err)
		// Don't fail the group creation
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Group created successfully",
		"group": fiber.Map{
			"id":          group.ID,
			"name":        group.Name,
			"description": group.Description,
			"wallet":      group.Wallet,
			"secret_key":  group.SecretKey,
			"status":      group.Status,
			"contract_id": contractID,
		},
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
	balance, err := services.CheckBalance(group.Wallet)
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

// GetGroupSecretKey returns the group's secret key (only for creators/admins)
func GetGroupSecretKey(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Check if user is admin/creator of the group
	var member models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?", 
		groupID, user.ID, []string{"creator", "admin"}).First(&member).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only group creators and admins can view the secret key",
		})
	}

	// Get group
	group, err := services.GetGroupByID(groupID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Group not found",
		})
	}

	return c.JSON(fiber.Map{
		"group_id":   group.ID,
		"wallet":     group.Wallet,
		"secret_key": group.SecretKey,
	})
}
