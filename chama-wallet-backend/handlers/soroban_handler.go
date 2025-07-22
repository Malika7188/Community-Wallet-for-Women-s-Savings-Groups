package handlers

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/services"
)

func ContributeHandler(c *fiber.Ctx) error {
	type RequestBody struct {
		UserAddress string `json:"user_address"`
		Amount      string `json:"amount"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Add your Soroban contract ID
	contractID := "CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4"

	args := []string{body.UserAddress, body.Amount}

	// Use CallSorobanFunctionWithAuth instead
	result, err := services.CallSorobanFunctionWithAuth(contractID, "contribute", body.UserAddress, args)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Contribution successful",
		"result":  result,
	})
}
