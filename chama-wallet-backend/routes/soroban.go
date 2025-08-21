package routes

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/handlers"
)

func SetupSorobanRoutes(app *fiber.App) {
	// Soroban contract interaction routes
	app.Post("/api/contribute", handlers.ContributeHandler)
	app.Post("/api/balance", handlers.BalanceHandler)
	app.Post("/api/withdraw", handlers.WithdrawHandler)
	app.Post("/api/history", handlers.HistoryHandler)
}
