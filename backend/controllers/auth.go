package controllers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"music-player-backend/database"
	"music-player-backend/middleware"
	"music-player-backend/models"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID int) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(middleware.GetJWTSecret()))
}

func GetUserIDFromContext(c *fiber.Ctx) (int, error) {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID")
	}
	return strconv.Atoi(userIDStr)
}

func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Username is required",
		})
	}

	if req.Email == "" {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Email is required",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Password must be at least 6 characters",
		})
	}

	var exists bool
	err := database.Pool.QueryRow(
		context.Background(),
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`,
		req.Username,
	).Scan(&exists)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to check username: " + err.Error(),
		})
	}

	if exists {
		return c.Status(409).JSON(models.Response{
			Success: false,
			Error:   "Username already exists",
		})
	}

	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`,
		req.Email,
	).Scan(&exists)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to check email: " + err.Error(),
		})
	}

	if exists {
		return c.Status(409).JSON(models.Response{
			Success: false,
			Error:   "Email already exists",
		})
	}

	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to hash password: " + err.Error(),
		})
	}

	var userID int
	err = database.Pool.QueryRow(
		context.Background(),
		`INSERT INTO users (username, email, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		req.Username, req.Email, passwordHash,
	).Scan(&userID)

	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to create user: " + err.Error(),
		})
	}

	token, err := GenerateToken(userID)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to generate token: " + err.Error(),
		})
	}

	user := models.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	return c.Status(201).JSON(models.Response{
		Success: true,
		Data: models.LoginResponse{
			Token: token,
			User:  user,
		},
	})
}

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(models.Response{
			Success: false,
			Error:   "Email and password are required",
		})
	}

	var user models.UserWithPassword
	err := database.Pool.QueryRow(
		context.Background(),
		`SELECT id, username, email, password_hash 
		 FROM users 
		 WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)

	if err != nil {
		return c.Status(401).JSON(models.Response{
			Success: false,
			Error:   "Invalid email or password",
		})
	}

	if !CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(401).JSON(models.Response{
			Success: false,
			Error:   "Invalid email or password",
		})
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			Success: false,
			Error:   "Failed to generate token: " + err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Data: models.LoginResponse{
			Token: token,
			User: models.User{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			},
		},
	})
}

func GetCurrentUser(c *fiber.Ctx) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(models.Response{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	var user models.User
	err = database.Pool.QueryRow(
		context.Background(),
		`SELECT id, username, email FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		return c.Status(404).JSON(models.Response{
			Success: false,
			Error:   "User not found",
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Data:    user,
	})
}
