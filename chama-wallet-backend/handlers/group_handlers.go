package handlers

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/services"
)

func ContributeToGroup(c *fiber.Ctx) error {
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

	err = services.SendPayment(payload.Secret, group.Wallet, payload.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Contribution successful",
		"group":   groupID,
		"from":    payload.From,
		"to":      group.Wallet,
		"amount":  payload.Amount,
	})
}
