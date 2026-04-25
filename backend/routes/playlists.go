package routes

import (
	"github.com/gofiber/fiber/v2"
	"music-player-backend/controllers"
)

func SetupPlaylistsRoutes(app *fiber.App) {
	playlists := app.Group("/api/playlists")

	playlists.Get("/my", controllers.GetMyPlaylists)
	playlists.Get("/popular", controllers.GetPopularPlaylists)
	
	playlists.Post("/", controllers.CreatePlaylist)
	playlists.Get("/:id", controllers.GetPlaylist)
	playlists.Get("/:id/songs", controllers.GetPlaylistSongs)
	playlists.Put("/:id", controllers.UpdatePlaylist)
	playlists.Delete("/:id", controllers.DeletePlaylist)
	
	playlists.Post("/:id/songs", controllers.AddSongToPlaylist)
	playlists.Delete("/:id/songs/:songId", controllers.RemoveSongFromPlaylist)
	playlists.Put("/:id/positions", controllers.UpdateSongPositions)
	
	playlists.Post("/:id/copy", controllers.CopyPlaylist)
}
