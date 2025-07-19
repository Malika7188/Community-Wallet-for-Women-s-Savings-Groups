package handlers

import (
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
		Name   string `json:"name"`
		Wallet string `json:"wallet"`
		Admin  string `json:"admin"` // Optional: If the creator is stored
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// ✅ Step 1: Deploy contract using CLI (or pre-deployed if needed)
	contractID, err := services.DeployChamaContract()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to deploy contract"})
	}

	// ✅ Step 2: Save group in DB with the new contractID
	group := models.Group{
		ID:         uuid.NewString(),
		Name:       payload.Name,
		Wallet:     payload.Wallet,
		ContractID: contractID,
	}

	if err := database.DB.Create(&group).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save group"})
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
	}

	if err := c.BodyParser(&body); err != nil || body.Wallet == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request. Wallet is required.",
		})
	}

	group, err := services.AddMemberToGroup(groupID, body.Wallet)
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
	balance, err := services.GetBalance(group.Wallet)
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
	groups, err := services.GetAllGroups()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch groups",
		})
	}

	return c.JSON(groups)
}
