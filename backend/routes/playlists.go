package routes

import (
	"github.com/gofiber/fiber/v2"
	"music-player-backend/controllers"
	"music-player-backend/middleware"
)

func SetupPlaylistsRoutes(app *fiber.App) {
	playlists := app.Group("/api/playlists")

	playlists.Get("/popular", controllers.GetPopularPlaylists)

	protectedPlaylists := playlists.Group("", middleware.JWTMiddleware())
	protectedPlaylists.Get("/my", controllers.GetMyPlaylists)
	protectedPlaylists.Post("/", controllers.CreatePlaylist)
	protectedPlaylists.Post("/:id/copy", controllers.CopyPlaylist)

	optionalPlaylists := playlists.Group("", middleware.OptionalJWTMiddleware())
	optionalPlaylists.Get("/:id", controllers.GetPlaylist)
	optionalPlaylists.Get("/:id/songs", controllers.GetPlaylistSongs)

	ownerProtectedPlaylists := playlists.Group("", middleware.JWTMiddleware())
	ownerProtectedPlaylists.Put("/:id", controllers.UpdatePlaylist)
	ownerProtectedPlaylists.Delete("/:id", controllers.DeletePlaylist)
	ownerProtectedPlaylists.Post("/:id/songs", controllers.AddSongToPlaylist)
	ownerProtectedPlaylists.Delete("/:id/songs/:songId", controllers.RemoveSongFromPlaylist)
	ownerProtectedPlaylists.Put("/:id/positions", controllers.UpdateSongPositions)
}
