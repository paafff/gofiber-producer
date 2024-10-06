package main

import (
	"fmt"
	"gofiber-producer/internal/auth"
	"gofiber-producer/internal/config"
	"gofiber-producer/internal/domain/services"
	"gofiber-producer/internal/handlers"
	"gofiber-producer/internal/infrastructure/database"
	"gofiber-producer/internal/infrastructure/router"
	"gofiber-producer/internal/repositories"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	config.LoadConfig()
	log.Println("Configuration loaded")

	// Inisialisasi database
	database.InitDatabase()
	log.Println("Database initialized")

	// Handle command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "seed":
			database.SeedUsers(database.DB, 10)
			fmt.Println("Seeded 10 users")
			return
		case "reset-db":
			database.ResetDatabase(database.DB)
			fmt.Println("Database reset.")
			return
		}
	}

	// Inisialisasi repository
	userRepository := repositories.NewUserRepository(database.DB)
	log.Println("User repository initialized")

	// Inisialisasi service
	userService := services.NewUserService(userRepository)
	authService := auth.NewAuthService(userRepository)
	log.Println("User service and Auth service initialized")

	// Inisialisasi handler
	userHandler := handlers.NewUserHandler(userService)
	authHandler := auth.NewAuthHandler(authService)
	log.Println("User handler and Auth handler initialized")

	// // Inisialisasi service
	// userService := services.NewUserService(userRepository)
	// log.Println("User service initialized")

	// // Inisialisasi handler
	// userHandler := handlers.NewUserHandler(userService)
	// log.Println("User handler initialized")

	// Inisialisasi Fiber
	app := fiber.New()
	log.Println("Fiber app initialized")

	// Setup routes
	router.SetupRoutes(app, userHandler, authHandler)
	log.Println("Routes set up")

	// Jalankan server
	log.Println("Server running on port 8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
