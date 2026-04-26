package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"music-player-backend/database"
	"music-player-backend/models"
)

func ListUsers(c *fiber.Ctx) error {
	rows, err := database.Pool.Query(
		context.Background(),
		`SELECT id, username FROM users ORDER BY id`,
	)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to list users: " + err.Error(),
		})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return c.Status(500).JSON(models.Response{
				Success: false,
				Error:   "Failed to scan user: " + err.Error(),
			})
		}
		users = append(users, user)
	}

	return c.JSON(models.Response{
		Success: true,
		Data:    users,
	})
}
