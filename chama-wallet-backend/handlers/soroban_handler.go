package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/services"
	"chama-wallet-backend/config"
)

// ContributeHandler handles direct Soroban contributions
func ContributeHandler(c *fiber.Ctx) error {
	type RequestBody struct {
		ContractID  string `json:"contract_id"`
		UserAddress string `json:"user_address"`
		Amount      string `json:"amount"`
		SecretKey   string `json:"secret_key,omitempty"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		fmt.Printf("❌ Failed to parse contribute request: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if body.ContractID == "" || body.UserAddress == "" || body.Amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: contract_id, user_address, and amount are required",
		})
	}

	// Validate amount
	amount, err := strconv.ParseFloat(body.Amount, 64)
	if err != nil || amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be a positive number",
		})
	}

	// Validate amount limits for mainnet
	if config.Config.IsMainnet {
		minAmount := 0.0000001 // Minimum XLM amount
		if amount < minAmount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Amount below minimum of %f XLM for mainnet", minAmount),
			})
		}
	}

	fmt.Printf("🔄 Processing direct Soroban contribution: %s XLM from %s to contract %s on %s\n",
		body.Amount, body.UserAddress, body.ContractID, config.Config.Network)
	args := []string{body.UserAddress, body.Amount}
	var result string

	// Use authenticated call if secret key provided, otherwise use regular call
	if body.SecretKey != "" {
		result, err = services.CallSorobanFunctionWithAuth(body.ContractID, "contribute", body.SecretKey, args)
	} else {
		result, err = services.CallSorobanFunction(body.ContractID, "contribute", args)
	}

	if err != nil {
		fmt.Printf("❌ Soroban contribution failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Blockchain transaction failed: %v", err),
		})
	}

	fmt.Printf("✅ Direct Soroban contribution successful on %s: %s\n", config.Config.Network, result)

	return c.JSON(fiber.Map{
		"message":     "Contribution successful",
		"contract_id": body.ContractID,
		"user":        body.UserAddress,
		"amount":      body.Amount,
		"tx_hash":     result,
		"network":     config.Config.Network,
		"timestamp":   fmt.Sprintf("%d", time.Now().Unix()),
	})
}

// BalanceHandler handles balance queries
func BalanceHandler(c *fiber.Ctx) error {
	type RequestBody struct {
		ContractID  string `json:"contract_id"`
		UserAddress string `json:"user_address"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if body.ContractID == "" || body.UserAddress == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: contract_id and user_address are required",
		})
	}

	result, err := services.GetBalance(body.ContractID, body.UserAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get balance: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"contract_id": body.ContractID,
		"user":        body.UserAddress,
		"balance":     result,
		"network":     config.Config.Network,
	})
}

// WithdrawHandler handles withdrawal requests
func WithdrawHandler(c *fiber.Ctx) error {
	type RequestBody struct {
		ContractID  string `json:"contract_id"`
		UserAddress string `json:"user_address"`
		Amount      string `json:"amount"`
		SecretKey   string `json:"secret_key"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if body.ContractID == "" || body.UserAddress == "" || body.Amount == "" || body.SecretKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: contract_id, user_address, amount, and secret_key are required",
		})
	}

	// Validate amount
	amount, err := strconv.ParseFloat(body.Amount, 64)
	if err != nil || amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be a positive number",
		})
	}

	// Additional validation for mainnet
	if config.Config.IsMainnet {
		minAmount := 0.0000001
		if amount < minAmount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Amount below minimum of %f XLM for mainnet", minAmount),
			})
		}
	}
	args := []string{body.UserAddress, body.Amount}
	result, err := services.CallSorobanFunctionWithAuth(body.ContractID, "withdraw", body.SecretKey, args)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Withdrawal failed: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Withdrawal successful",
		"contract_id": body.ContractID,
		"user":        body.UserAddress,
		"amount":      body.Amount,
		"new_balance": result,
		"network":     config.Config.Network,
	})
}

// HistoryHandler handles contribution history requests
func HistoryHandler(c *fiber.Ctx) error {
	type RequestBody struct {
		ContractID  string `json:"contract_id"`
		UserAddress string `json:"user_address"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if body.ContractID == "" || body.UserAddress == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: contract_id and user_address are required",
		})
	}

	result, err := services.GetContributionHistory(body.ContractID, body.UserAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get history: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"contract_id": body.ContractID,
		"user":        body.UserAddress,
		"history":     result,
		"network":     config.Config.Network,
	})
}
