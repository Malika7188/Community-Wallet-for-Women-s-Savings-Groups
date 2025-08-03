package routes

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/handlers"
	"chama-wallet-backend/middleware"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 Community Wallet API is running.")
	})
	
	// Public wallet routes
	app.Post("/create-wallet", handlers.CreateWallet)
	app.Get("/balance/:address", middleware.OptionalAuthMiddleware(), handlers.GetBalance)
	app.Get("/generate-keypair", handlers.GenerateKeypair)
	app.Post("/fund/:address", middleware.OptionalAuthMiddleware(), handlers.FundAccount)
	app.Get("/transactions/:address", middleware.OptionalAuthMiddleware(), handlers.GetTransactionHistory)
	
	// Protected wallet routes
	app.Post("/transfer", middleware.AuthMiddleware(), handlers.TransferFunds)
	app.Post("/api/contribute", handlers.ContributeHandler)

}
