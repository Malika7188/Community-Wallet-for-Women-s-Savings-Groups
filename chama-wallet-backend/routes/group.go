// routes/group.go
package routes

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/handlers"
)

func GroupRoutes(app *fiber.App) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	app.Post("/group/create", handlers.CreateGroup)
	app.Post("/group/:id/contribute", handlers.ContributeToGroup)
	app.Post("/group/:id/join", handlers.AddMember)
	app.Get("/group/:id/balance", handlers.GetGroupBalance)
	app.Get("/groups", handlers.GetAllGroups)

}
