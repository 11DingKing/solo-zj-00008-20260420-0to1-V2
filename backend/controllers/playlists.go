package controllers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"music-player-backend/database"
	"music-player-backend/models"
)

const currentUserID = 1

func CreatePlaylist(c *fiber.Ctx) error {
	var req models.CreatePlaylistRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
	}

	var playlistID int
	err := database.Pool.QueryRow(
		context.Background(),
		`INSERT INTO playlists (name, description, is_public, owner_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		req.Name, req.Description, req.IsPublic, currentUserID,
	).Scan(&playlistID)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to create playlist: " + err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Data:    map[string]int{"id": playlistID},
	})
}

func GetPlaylist(c *fiber.Ctx) error {
	id := c.Params("id")
	playlistID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid playlist ID",
		})
	}

	var playlist models.Playlist
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT id, name, description, is_public, owner_id
		 FROM playlists WHERE id = $1`,
		playlistID,
	).Scan(&playlist.ID, &playlist.Name, &playlist.Description, &playlist.IsPublic, &playlist.OwnerID)

	if err != nil {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Playlist not found",
		})
	}

	if !playlist.IsPublic && playlist.OwnerID != currentUserID {
		return c.Status(403).JSON(models.Response{
			Success: false,
			Error:   "Access denied",
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Data:    playlist,
	})
}

func GetPlaylistSongs(c *fiber.Ctx) error {
	id := c.Params("id")
	playlistID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid playlist ID",
		})
	}

	var playlist models.Playlist
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT is_public, owner_id FROM playlists WHERE id = $1`,
		playlistID,
	).Scan(&playlist.IsPublic, &playlist.OwnerID)

	if err != nil {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Playlist not found",
		})
	}

	if !playlist.IsPublic && playlist.OwnerID != currentUserID {
		return c.Status(403).JSON(models.Response{
			Success: false,
			Error:   "Access denied",
		})
	}

	rows, err := database.Pool.Query(
		context.Background(),
		`SELECT s.id, s.name, s.artist, s.duration, ps.position
		 FROM playlist_songs ps
		 JOIN songs s ON ps.song_id = s.id
		 WHERE ps.playlist_id = $1
		 ORDER BY ps.position`,
		playlistID,
	)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to get playlist songs: " + err.Error(),
		})
	}
	defer rows.Close()

	var songs []models.PlaylistSongDetail
	var totalDuration int
	for rows.Next() {
		var song models.PlaylistSongDetail
		if err := rows.Scan(&song.SongID, &song.Name, &song.Artist, &song.Duration, &song.Position); err != nil {
			return c.Status(500).JSON(models.Response{
				Success: false,
				Error:   "Failed to scan song: " + err.Error(),
			})
		}
		songs = append(songs, song)
		totalDuration += song.Duration
	}

	return c.JSON(models.Response{
		Success: true,
		Data: map[string]interface{}{
			"songs":          songs,
			"total_duration": totalDuration,
		},
	})
}

func UpdatePlaylist(c *fiber.Ctx) error {
	id := c.Params("id")
	playlistID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid playlist ID",
		})
	}

	var ownerID int
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT owner_id FROM playlists WHERE id = $1`,
		playlistID,
	).Scan(&ownerID)

	if err != nil {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Playlist not found",
		})
	}

	if ownerID != currentUserID {
		return c.Status(403).JSON(models.Response{
			Success: false,
			Error:   "Access denied",
		})
	}

	var req models.CreatePlaylistRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
	}

	commandTag, err := database.Pool.Exec(
		context.Background(),
		`UPDATE playlists 
		 SET name = $1, description = $2, is_public = $3
		 WHERE id = $4`,
		req.Name, req.Description, req.IsPublic, playlistID,
	)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to update playlist: " + err.Error(),
		})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Playlist not found",
		})
	}

	return c.JSON(models.Response{
		Success: true,
	})
}

func DeletePlaylist(c *fiber.Ctx) error {
	id := c.Params("id")
	playlistID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid playlist ID",
		})
	}

	var ownerID int
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT owner_id FROM playlists WHERE id = $1`,
		playlistID,
	).Scan(&ownerID)

	if err != nil {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Playlist not found",
		})
	}

	if ownerID != currentUserID {
		return c.Status(403).JSON(models.Response{
			Success: false,
			Error:   "Access denied",
		})
	}

	commandTag, err := database.Pool.Exec(
		context.Background(),
		`DELETE FROM playlists WHERE id = $1`,
		playlistID,
	)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to delete playlist: " + err.Error(),
		})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Playlist not found",
		})
	}

	return c.JSON(models.Response{
		Success: true,
	})
}

func AddSongToPlaylist(c *fiber.Ctx) error {
	id := c.Params("id")
	playlistID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid playlist ID",
		})
	}

	var ownerID int
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT owner_id FROM playlists WHERE id = $1`,
		playlistID,
	).Scan(&ownerID)

	if err != nil {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Playlist not found",
		})
	}

	if ownerID != currentUserID {
		return c.Status(403).JSON(models.Response{
			Success: false,
			Error:   "Access denied",
		})
	}

	var req map[string]int
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
	}

	songID, ok := req["song_id"]
	if !ok {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "song_id is required",
		})
	}

	var exists bool
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT EXISTS(SELECT 1 FROM playlist_songs WHERE playlist_id = $1 AND song_id = $2)`,
		playlistID, songID,
	).Scan(&exists)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to check song existence: " + err.Error(),
		})
	}

	if exists {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Song already in playlist",
		})
	}

	var maxPosition int
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT COALESCE(MAX(position), -1) FROM playlist_songs WHERE playlist_id = $1`,
		playlistID,
	).Scan(&maxPosition)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to get max position: " + err.Error(),
		})
	}

	_, err = database.Pool.Exec(
		context.Background(),
		`INSERT INTO playlist_songs (playlist_id, song_id, position)
		 VALUES ($1, $2, $3)`,
		playlistID, songID, maxPosition+1,
	)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to add song to playlist: " + err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
	})
}

func RemoveSongFromPlaylist(c *fiber.Ctx) error {
	id := c.Params("id")
	playlistID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid playlist ID",
		})
	}

	songIDStr := c.Params("songId")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid song ID",
		})
	}

	var ownerID int
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT owner_id FROM playlists WHERE id = $1`,
		playlistID,
	).Scan(&ownerID)

	if err != nil {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Playlist not found",
		})
	}

	if ownerID != currentUserID {
		return c.Status(403).JSON(models.Response{
			Success: false,
			Error:   "Access denied",
		})
	}

	var position int
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT position FROM playlist_songs WHERE playlist_id = $1 AND song_id = $2`,
		playlistID, songID,
	).Scan(&position)

	if err != nil {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Song not found in playlist",
		})
	}

	_, err = database.Pool.Exec(
		context.Background(),
		`DELETE FROM playlist_songs WHERE playlist_id = $1 AND song_id = $2`,
		playlistID, songID,
	)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to remove song from playlist: " + err.Error(),
		})
	}

	_, err = database.Pool.Exec(
		context.Background(),
		`UPDATE playlist_songs 
		 SET position = position - 1 
		 WHERE playlist_id = $1 AND position > $2`,
		playlistID, position,
	)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to update positions: " + err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
	})
}

func UpdateSongPositions(c *fiber.Ctx) error {
	id := c.Params("id")
	playlistID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid playlist ID",
		})
	}

	var ownerID int
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT owner_id FROM playlists WHERE id = $1`,
		playlistID,
	).Scan(&ownerID)

	if err != nil {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Playlist not found",
		})
	}

	if ownerID != currentUserID {
		return c.Status(403).JSON(models.Response{
			Success: false,
			Error:   "Access denied",
		})
	}

	var req []models.PlaylistSong
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
	}

	tx, err := database.Pool.Begin(context.Background())
	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to start transaction: " + err.Error(),
		})
	}
	defer tx.Rollback(context.Background())

	for _, item := range req {
		_, err = tx.Exec(
			context.Background(),
			`UPDATE playlist_songs 
			 SET position = $1 
			 WHERE playlist_id = $2 AND song_id = $3`,
			item.Position, playlistID, item.SongID,
		)
		if err != nil {
			return c.Status(500).JSON(models.Response{
				Success: false,
				Error:   "Failed to update position: " + err.Error(),
			})
		}
	}

	if err = tx.Commit(context.Background()); err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to commit transaction: " + err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
	})
}

func CopyPlaylist(c *fiber.Ctx) error {
	id := c.Params("id")
	playlistID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid playlist ID",
		})
	}

	var playlist models.Playlist
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT name, description, is_public, owner_id 
		 FROM playlists WHERE id = $1`,
		playlistID,
	).Scan(&playlist.Name, &playlist.Description, &playlist.IsPublic, &playlist.OwnerID)

	if err != nil {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Playlist not found",
		})
	}

	if !playlist.IsPublic && playlist.OwnerID != currentUserID {
		return c.Status(403).JSON(models.Response{
			Success: false,
			Error:   "Access denied - playlist is not public",
		})
	}

	if playlist.OwnerID == currentUserID {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Cannot copy your own playlist",
		})
	}

	tx, err := database.Pool.Begin(context.Background())
	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to start transaction: " + err.Error(),
		})
	}
	defer tx.Rollback(context.Background())

	var newPlaylistID int
	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO playlists (name, description, is_public, owner_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		playlist.Name+" (Copy)", playlist.Description, false, currentUserID,
	).Scan(&newPlaylistID)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to create playlist copy: " + err.Error(),
		})
	}

	songs, err := tx.Query(
		context.Background(),
		`SELECT song_id, position 
		 FROM playlist_songs 
		 WHERE playlist_id = $1
		 ORDER BY position`,
		playlistID,
	)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to get playlist songs: " + err.Error(),
		})
	}
	defer songs.Close()

	for songs.Next() {
		var songID int
		var position int
		if err := songs.Scan(&songID, &position); err != nil {
			return c.Status(500).JSON(models.Response{
				Success: false,
				Error:   "Failed to scan song: " + err.Error(),
			})
		}

		_, err = tx.Exec(
			context.Background(),
			`INSERT INTO playlist_songs (playlist_id, song_id, position)
			 VALUES ($1, $2, $3)`,
			newPlaylistID, songID, position,
		)
		if err != nil {
			return c.Status(500).JSON(models.Response{
				Success: false,
				Error:   "Failed to copy song: " + err.Error(),
			})
		}
	}

	if err = tx.Commit(context.Background()); err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to commit transaction: " + err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Data:    map[string]int{"id": newPlaylistID},
	})
}

func GetMyPlaylists(c *fiber.Ctx) error {
	rows, err := database.Pool.Query(
		context.Background(),
		`SELECT p.id, p.name, p.description, p.is_public, p.owner_id,
		       COUNT(ps.playlist_id) as song_count
		 FROM playlists p
		 LEFT JOIN playlist_songs ps ON p.id = ps.playlist_id
		 WHERE p.owner_id = $1
		 GROUP BY p.id
		 ORDER BY p.id DESC`,
		currentUserID,
	)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to get my playlists: " + err.Error(),
		})
	}
	defer rows.Close()

	var playlists []models.Playlist
	for rows.Next() {
		var playlist models.Playlist
		if err := rows.Scan(
			&playlist.ID, &playlist.Name, &playlist.Description, 
			&playlist.IsPublic, &playlist.OwnerID, &playlist.SongCount,
		); err != nil {
			return c.Status(500).JSON(models.Response{
				Success: false,
				Error:   "Failed to scan playlist: " + err.Error(),
			})
		}
		playlists = append(playlists, playlist)
	}

	return c.JSON(models.Response{
		Success: true,
		Data:    playlists,
	})
}

func GetPopularPlaylists(c *fiber.Ctx) error {
	rows, err := database.Pool.Query(
		context.Background(),
		`SELECT p.id, p.name, p.description, p.is_public, p.owner_id,
		       COUNT(ps.playlist_id) as song_count
		 FROM playlists p
		 LEFT JOIN playlist_songs ps ON p.id = ps.playlist_id
		 WHERE p.is_public = true
		 GROUP BY p.id
		 ORDER BY song_count DESC
		 LIMIT 10`,
	)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to get popular playlists: " + err.Error(),
		})
	}
	defer rows.Close()

	var playlists []models.Playlist
	for rows.Next() {
		var playlist models.Playlist
		if err := rows.Scan(
			&playlist.ID, &playlist.Name, &playlist.Description, 
			&playlist.IsPublic, &playlist.OwnerID, &playlist.SongCount,
		); err != nil {
			return c.Status(500).JSON(models.Response{
				Success: false,
				Error:   "Failed to scan playlist: " + err.Error(),
			})
		}
		playlists = append(playlists, playlist)
	}

	return c.JSON(models.Response{
		Success: true,
		Data:    playlists,
	})
}
