package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gofiber-producer/internal/infrastructure/database"

	"github.com/gofiber/fiber/v2"
)

type SessionData struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

func SessionMiddleware(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	sessionJSON, err := database.RedisDB.Get(context.Background(), sessionID).Result()
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var sessionData SessionData
	err = json.Unmarshal([]byte(sessionJSON), &sessionData)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Set session data in context
	c.Locals("sessionData", sessionData)

	return c.Next()
}

func CreateSession(sessionID string, data SessionData, expiration time.Duration) error {
	sessionJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return database.RedisDB.Set(context.Background(), sessionID, sessionJSON, expiration).Err()
}

func DeleteSession(sessionID string) error {
	return database.RedisDB.Del(context.Background(), sessionID).Err()
}
