package routes

import (
	"github.com/gofiber/fiber/v2"
	"music-player-backend/controllers"
)

func SetupUsersRoutes(app *fiber.App) {
	users := app.Group("/api/users")

	users.Get("/me", controllers.GetCurrentUser)
	users.Get("/", controllers.ListUsers)
}
