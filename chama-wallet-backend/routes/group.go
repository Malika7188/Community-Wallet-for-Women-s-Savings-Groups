// routes/group.go
package routes

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/handlers"
	"chama-wallet-backend/middleware"
)

func GroupRoutes(app *fiber.App) {

	// Public routes (can be accessed without authentication)
	app.Get("/groups", middleware.OptionalAuthMiddleware(), handlers.GetAllGroups)
	app.Get("/group/:id/balance", middleware.OptionalAuthMiddleware(), handlers.GetGroupBalance)

	// Protected routes (require authentication)
	app.Post("/group/create", middleware.AuthMiddleware(), handlers.CreateGroup)
	app.Post("/group/:id/contribute", middleware.AuthMiddleware(), handlers.ContributeToGroup)
	app.Post("/group/:id/join", middleware.AuthMiddleware(), handlers.AddMember)


}
