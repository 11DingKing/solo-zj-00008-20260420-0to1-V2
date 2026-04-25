package routes

import (
	"github.com/gofiber/fiber/v2"
	"music-player-backend/controllers"
	"music-player-backend/middleware"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/api/auth")

	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)

	protectedAuth := auth.Group("", middleware.JWTMiddleware())
	protectedAuth.Get("/me", controllers.GetCurrentUser)
}
