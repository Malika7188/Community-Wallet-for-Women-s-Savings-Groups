// routes/group.go
package routes

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/handlers"
	"chama-wallet-backend/middleware"
)

func GroupRoutes(app *fiber.App) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// Public routes (can be accessed without authentication)
	app.Get("/groups", middleware.AuthMiddleware(), handlers.GetAllGroups)
	app.Get("/group/:id", middleware.OptionalAuthMiddleware(), handlers.GetGroupDetails)
	app.Get("/group/:id/balance", middleware.OptionalAuthMiddleware(), handlers.GetGroupBalance)

	// Protected routes (require authentication)
	app.Post("/group/create", middleware.AuthMiddleware(), handlers.CreateGroup)
	app.Get("/user/groups", middleware.AuthMiddleware(), handlers.GetUserGroups)
	app.Post("/group/:id/contribute", middleware.AuthMiddleware(), handlers.ContributeToGroup)
	app.Post("/group/:id/join", middleware.AuthMiddleware(), handlers.JoinGroup)

	// New routes
	app.Post("/group/:id/invite", middleware.AuthMiddleware(), handlers.InviteToGroup)
	app.Get("/group/:id/non-members", middleware.AuthMiddleware(), handlers.GetNonGroupMembers)
	app.Post("/group/:id/approve", middleware.AuthMiddleware(), handlers.ApproveGroup)
	app.Post("/group/:id/activate", middleware.AuthMiddleware(), handlers.ActivateGroup)
	app.Post("/group/:id/nominate-admin", middleware.AuthMiddleware(), handlers.NominateAdmin)
	app.Post("/group/:id/approve-member", middleware.AuthMiddleware(), handlers.ApproveMember)
	app.Post("/group/:id/payout-request", middleware.AuthMiddleware(), handlers.CreatePayoutRequest)
	app.Post("/payout/:id/approve", middleware.AuthMiddleware(), handlers.ApprovePayoutRequest)
	app.Get("/group/:id/payout-requests", middleware.AuthMiddleware(), handlers.GetPayoutRequests)
	app.Get("/group/:id/payout-schedule", middleware.AuthMiddleware(), handlers.GetPayoutSchedule)

	// Notification routes
	app.Get("/notifications", middleware.AuthMiddleware(), handlers.GetNotifications)
	app.Put("/notifications/:id/read", middleware.AuthMiddleware(), handlers.MarkNotificationRead)
	app.Delete("/notifications/:id", middleware.AuthMiddleware(), handlers.DeleteNotification)
	app.Get("/invitations", middleware.AuthMiddleware(), handlers.GetUserInvitations)
	app.Post("/invitations/:id/accept", middleware.AuthMiddleware(), handlers.AcceptInvitation)
	app.Post("/invitations/:id/reject", middleware.AuthMiddleware(), handlers.RejectInvitation)

	// Contribution round routes
	app.Post("/group/:id/contribute-round", middleware.AuthMiddleware(), handlers.ContributeToRound)
	app.Get("/group/:id/round-status", middleware.AuthMiddleware(), handlers.GetRoundStatus)
	app.Post("/group/:id/authorize-payout", middleware.AuthMiddleware(), handlers.AuthorizeRoundPayout)

	// Add this route for group secret key access
	app.Get("/group/:id/secret", middleware.AuthMiddleware(), handlers.GetGroupSecretKey)
}
