package database

import (
	"log"

	"gofiber-producer/internal/domain/models"

	"github.com/bxcodec/faker/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB, count int) {
	password := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Could not hash password: %v", err)
	}

	for i := 0; i < count; i++ {
		user := models.User{
			Name:     faker.Name(),
			Email:    faker.Email(),
			Password: string(hashedPassword),
		}
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Could not seed user: %v", err)
		}
	}
	log.Printf("Seeded %d users", count)
}

// Fungsi untuk mereset database
func ResetDatabase(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.User{})
	if err != nil {
		log.Fatalf("Could not drop table: %v", err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Could not migrate table: %v", err)
	}
	log.Println("Database reset.")
}
