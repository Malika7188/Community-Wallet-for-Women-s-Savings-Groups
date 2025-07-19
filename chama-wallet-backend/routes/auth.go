package routes

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/handlers"
	"chama-wallet-backend/middleware"
)

// AuthRoutes sets up authentication routes
func AuthRoutes(app *fiber.App) {
	// Public routes
	app.Post("/auth/register", handlers.Register)
	app.Post("/auth/login", handlers.Login)
	app.Post("/auth/logout", handlers.Logout)

	// Protected routes
	auth := app.Group("/auth", middleware.AuthMiddleware())
	auth.Get("/profile", handlers.GetProfile)
	auth.Put("/profile", handlers.UpdateProfile)
}