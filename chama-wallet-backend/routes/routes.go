package routes

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/handlers"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 Community Wallet API is running.")
	})
	app.Post("/create-wallet", handlers.CreateWallet)
	app.Get("/balance/:address", handlers.GetBalance)
	app.Post("/transfer", handlers.TransferFunds)
	app.Get("/generate-keypair", handlers.GenerateKeypair)
	app.Post("/fund/:address", handlers.FundAccount)
	app.Get("/transactions/:address", handlers.GetTransactionHistory)
	app.Post("/api/contribute", handlers.ContributeHandler)

}
