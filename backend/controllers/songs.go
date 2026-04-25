package controllers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"music-player-backend/database"
	"music-player-backend/models"
)

func CreateSong(c *fiber.Ctx) error {
	var req models.CreateSongRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
	}

	var songID int
	err := database.Pool.QueryRow(
		context.Background(),
		`INSERT INTO songs (name, artist, album, duration, cover_url, audio_file_url)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id`,
		req.Name, req.Artist, req.Album, req.Duration, req.CoverURL, req.AudioFileURL,
	).Scan(&songID)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to create song: " + err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Data:    map[string]int{"id": songID},
	})
}

func GetSong(c *fiber.Ctx) error {
	id := c.Params("id")
	songID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid song ID",
		})
	}

	var song models.Song
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT id, name, artist, album, duration, cover_url, audio_file_url
		 FROM songs WHERE id = $1`,
		songID,
	).Scan(&song.ID, &song.Name, &song.Artist, &song.Album, &song.Duration, &song.CoverURL, &song.AudioFileURL)

	if err != nil {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Song not found",
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Data:    song,
	})
}

func UpdateSong(c *fiber.Ctx) error {
	id := c.Params("id")
	songID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid song ID",
		})
	}

	var req models.CreateSongRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
	}

	commandTag, err := database.Pool.Exec(
		context.Background(),
		`UPDATE songs 
		 SET name = $1, artist = $2, album = $3, duration = $4, cover_url = $5, audio_file_url = $6
		 WHERE id = $7`,
		req.Name, req.Artist, req.Album, req.Duration, req.CoverURL, req.AudioFileURL, songID,
	)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to update song: " + err.Error(),
		})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Song not found",
		})
	}

	return c.JSON(models.Response{
		Success: true,
	})
}

func DeleteSong(c *fiber.Ctx) error {
	id := c.Params("id")
	songID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid song ID",
		})
	}

	commandTag, err := database.Pool.Exec(
		context.Background(),
		`DELETE FROM songs WHERE id = $1`,
		songID,
	)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to delete song: " + err.Error(),
		})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "Song not found",
		})
	}

	return c.JSON(models.Response{
		Success: true,
	})
}

func ListSongs(c *fiber.Ctx) error {
	search := c.Query("search")
	
	var rows []models.Song
	var err error

	if search != "" {
		searchPattern := "%" + search + "%"
		rowsData, queryErr := database.Pool.Query(
			context.Background(),
			`SELECT id, name, artist, album, duration, cover_url, audio_file_url
			 FROM songs 
			 WHERE name ILIKE $1 OR artist ILIKE $2
			 ORDER BY artist, name`,
			searchPattern, searchPattern,
		)
		err = queryErr
		if err == nil {
			defer rowsData.Close()
			for rowsData.Next() {
				var song models.Song
				if scanErr := rowsData.Scan(&song.ID, &song.Name, &song.Artist, &song.Album, &song.Duration, &song.CoverURL, &song.AudioFileURL); scanErr != nil {
					err = scanErr
					break
				}
				rows = append(rows, song)
			}
		}
	} else {
		rowsData, queryErr := database.Pool.Query(
			context.Background(),
			`SELECT id, name, artist, album, duration, cover_url, audio_file_url
			 FROM songs 
			 ORDER BY artist, name`,
		)
		err = queryErr
		if err == nil {
			defer rowsData.Close()
			for rowsData.Next() {
				var song models.Song
				if scanErr := rowsData.Scan(&song.ID, &song.Name, &song.Artist, &song.Album, &song.Duration, &song.CoverURL, &song.AudioFileURL); scanErr != nil {
					err = scanErr
					break
				}
				rows = append(rows, song)
			}
		}
	}

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to list songs: " + err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Data:    rows,
	})
}
