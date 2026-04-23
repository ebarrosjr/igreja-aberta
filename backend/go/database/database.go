package database

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize(dbPath string) error {
	if dbPath == "" {
		dbPath = "igreja_aberta.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	DB = db
	return nil
}

func MigrateAll() error {
	return DB.AutoMigrate(&Congregation{}, &User{}, &RefreshToken{}, &PasswordResetToken{})
}

func SeedInitialData() error {
	const congregationEmail = "admin@igrejaaberta.local"
	const adminEmail = "ebarrosjr@gmail.com"

	congregation, err := GetCongregationByEmail(congregationEmail)
	if err != nil {
		return err
	}

	if congregation == nil {
		congregation, err = CreateCongregation(
			"Igreja Aberta",
			"Congregacao inicial da aplicacao",
			"",
			congregationEmail,
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		)
		if err != nil {
			return err
		}
	}

	adminExists, err := UserExistsByEmail(adminEmail)
	if err != nil {
		return err
	}

	if adminExists {
		return nil
	}

	adminPassword := os.Getenv("SEED_ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = os.Getenv("Seed__AdminPassword")
	}
	if adminPassword == "" {
		adminPassword = "Admin@123456"
	}

	hashedPassword, err := HashPassword(adminPassword)
	if err != nil {
		return err
	}

	_, err = CreateUser(congregation.Id, "Administrador", adminEmail, hashedPassword, "active")
	return err
}
