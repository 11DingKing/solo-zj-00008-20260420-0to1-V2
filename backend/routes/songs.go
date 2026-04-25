package routes

import (
	"github.com/gofiber/fiber/v2"
	"music-player-backend/controllers"
	"music-player-backend/middleware"
)

func SetupSongsRoutes(app *fiber.App) {
	songs := app.Group("/api/songs")

	songs.Get("/", controllers.ListSongs)
	songs.Get("/:id", controllers.GetSong)

	protectedSongs := songs.Group("", middleware.JWTMiddleware())
	protectedSongs.Post("/", controllers.CreateSong)
	protectedSongs.Put("/:id", controllers.UpdateSong)
	protectedSongs.Delete("/:id", controllers.DeleteSong)
}
