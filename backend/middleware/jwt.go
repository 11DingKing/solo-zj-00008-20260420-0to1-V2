package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"music-player-backend/models"
)

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}
	return secret
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(models.Response{
				Success: false,
				Error:   "Unauthorized - missing token",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(models.Response{
				Success: false,
				Error:   "Unauthorized - invalid token format",
			})
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(
			tokenString,
			&jwt.RegisteredClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(GetJWTSecret()), nil
			},
		)

		if err != nil {
			return c.Status(401).JSON(models.Response{
				Success: false,
				Error:   "Unauthorized - invalid token",
			})
		}

		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
			userIDStr := claims.Subject
			c.Locals("user_id", userIDStr)
			return c.Next()
		}

		return c.Status(401).JSON(models.Response{
			Success: false,
			Error:   "Unauthorized - invalid token claims",
		})
	}
}

func OptionalJWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(
			tokenString,
			&jwt.RegisteredClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(GetJWTSecret()), nil
			},
		)

		if err != nil {
			return c.Next()
		}

		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
			userIDStr := claims.Subject
			c.Locals("user_id", userIDStr)
		}

		return c.Next()
	}
}
