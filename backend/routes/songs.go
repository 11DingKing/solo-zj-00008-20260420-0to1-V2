package routes

import (
	"github.com/gofiber/fiber/v2"
	"music-player-backend/controllers"
)

func SetupSongsRoutes(app *fiber.App) {
	songs := app.Group("/api/songs")

	songs.Get("/", controllers.ListSongs)
	songs.Post("/", controllers.CreateSong)
	songs.Get("/:id", controllers.GetSong)
	songs.Put("/:id", controllers.UpdateSong)
	songs.Delete("/:id", controllers.DeleteSong)
}
