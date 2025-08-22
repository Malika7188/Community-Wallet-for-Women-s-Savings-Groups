package routes

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/handlers"
	"chama-wallet-backend/middleware"
)

// AutomatedPayoutRoutes sets up routes for automated payout functionality
func AutomatedPayoutRoutes(app *fiber.App) {
	// Protected routes - require authentication
	payout := app.Group("/payout", middleware.AuthMiddleware())
	
	// Get pending automatic payouts for admin approval
	payout.Get("/pending", handlers.GetPendingAutomaticPayouts)
	
	// Approve and execute automatic payout (single admin approval)
	payout.Post("/:id/approve-automatic", handlers.ApproveAutomaticPayout)
	
	// Get detailed payout status
	payout.Get("/:id/status", handlers.GetPayoutStatus)
	
	// Manually trigger automatic payout check (admin only)
	payout.Post("/trigger-check", handlers.TriggerAutomaticPayoutCheck)
}