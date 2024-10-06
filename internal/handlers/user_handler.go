package handlers

import (
	"gofiber-producer/internal/domain/models"
	"gofiber-producer/internal/domain/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	return c.JSON(user)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not hash password")
	}
	user.Password = string(hashedPassword)

	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not create user")
	}

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}
	user.ID = uint(id)

	updatedUser, err := h.userService.UpdateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not update user")
	}

	return c.JSON(updatedUser)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not delete user")
	}

	return c.JSON(fiber.Map{"message": "user deleted successfully"})
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not retrieve users")
	}

	return c.JSON(users)
}

// func ( h *UserHandler) Login(c *fiber.Ctx) error {
// 	login := new(models.Login)
// 	if err := c.BodyParser(login); err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
// 	}

// 	user, err := h.userService.GetUserByEmail(login.Email)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).SendString("User not found")
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
// 		return c.Status(fiber.StatusUnauthorized).SendString("Invalid password")
// 	}

// 	return c.JSON(user)
// }
