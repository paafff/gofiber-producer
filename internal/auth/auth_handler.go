package auth

import (
	"encoding/json"
	"gofiber-producer/internal/domain/models"
	"gofiber-producer/internal/infrastructure/database"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func generateSessionID() string {
	return uuid.New().String()
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	user, err := h.authService.Authenticate(credentials.Email, credentials.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// Generate token (JWT) and return it
	token, err := h.authService.GenerateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
	}

	// Create session
	sessionID := generateSessionID() // Implement this function to generate a unique session ID
	sessionData := SessionData{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
	}
	err = CreateSession(sessionID, sessionData, 24*7*time.Hour) // Set session expiration to 24 hours
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create session"})
	}

	// Set session ID as cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HTTPOnly: true,
	})

	// Save login activity to MongoDB
	activity := models.Activity{
		UserID:    user.ID,
		Email:     user.Email,
		LoginTime: time.Now(),
	}

	err = database.SaveLoginActivity(activity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not save login activity"})
	}

	// Log the successful save
	log.Printf("Login activity saved: UserID=%d, Email=%s, LoginTime=%s", activity.UserID, activity.Email, activity.LoginTime)

	return c.JSON(fiber.Map{"token": token, "user": user})
}

// func (h *AuthHandler) Login(c *fiber.Ctx) error {
// 	var credentials struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	if err := c.BodyParser(&credentials); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
// 	}

// 	user, err := h.authService.Authenticate(credentials.Email, credentials.Password)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
// 	}

// 	// Generate token (JWT) and return it
// 	token, err := h.authService.GenerateJWT(user)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
// 	}

// 	// Generate token (JWT, etc.) and return it
// 	// token := "dummy-token" // Replace with actual token generation logic
// 	return c.JSON(fiber.Map{"token": token, "user": user})
// }

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	createdUser, err := h.authService.Register(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
	}

	// Serialize user data to JSON
	userData, err := json.Marshal(createdUser)
	if err != nil {
		log.Printf("Failed to serialize user data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not serialize user data"})
	}

	// Publish message to RabbitMQ
	err = database.PublishMessage(string(userData))
	if err != nil {
		log.Printf("Failed to publish message to RabbitMQ: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not publish message"})
	}

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

// func (h *AuthHandler) Register(c *fiber.Ctx) error {
// 	user := new(models.User)
// 	if err := c.BodyParser(user); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
// 	}

// 	createdUser, err := h.authService.Register(user)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(createdUser)
// }
