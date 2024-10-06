package router

import (
	"gofiber-producer/internal/auth"
	"gofiber-producer/internal/handlers"
	"gofiber-producer/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler, authHandler *auth.AuthHandler) {
	api := app.Group("/api")
	userGroup := app.Group("/")

	// Public routes
	api.Post("/login", authHandler.Login)
	api.Post("/register", authHandler.Register)

	// Protected routes with token
	// api.Use(auth.AuthMiddleware)
	// routes.UserRoutes(api, userHandler)

	// Session routes
	userGroup.Use(auth.SessionMiddleware)
	routes.UserRoutes(userGroup, userHandler)
}
