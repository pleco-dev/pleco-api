package config

import (
	"go-auth-app/models"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	var user models.User

	DB.Where("email = ?", "admin@mail.com").First(&user)
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if user.ID == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

		admin := models.User{
			Name:     "Super Admin",
			Email:    email,
			Password: string(hashedPassword),
			Role:     "admin",
		}

		DB.Create(&admin)
	}
}
