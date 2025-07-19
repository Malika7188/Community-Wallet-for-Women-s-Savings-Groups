package routes

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/services"
)

func SetupSorobanRoutes(app *fiber.App) {
	contractID := "CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4"

	app.Post("/contribute", func(c *fiber.Ctx) error {
		var body struct {
			Amount  string `json:"amount"`
			Address string `json:"address"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
		}

		resp, err := services.CallSorobanFunction(contractID, "contribute", []string{body.Address, body.Amount})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"result": resp})
	})

	app.Get("/balance/:address", func(c *fiber.Ctx) error {
		address := c.Params("address")
		resp, err := services.CallSorobanFunction(contractID, "get_balance", []string{address})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"balance": resp})
	})

	app.Post("/withdraw", func(c *fiber.Ctx) error {
		var body struct {
			Amount  string `json:"amount"`
			Address string `json:"address"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
		}
		resp, err := services.CallSorobanFunction(contractID, "withdraw", []string{body.Address, body.Amount})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"result": resp})
	})

	app.Get("/history/:address", func(c *fiber.Ctx) error {
		address := c.Params("address")
		resp, err := services.CallSorobanFunction(contractID, "get_contribution_history", []string{address})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"history": resp})
	})
}
