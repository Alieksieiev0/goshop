package database

import (
	"fmt"
	"os"

	"github.com/Alieksieiev0/goshop/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Setup(db *gorm.DB) error {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error; err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Category{}, models.Product{}); err != nil {
		return err
	}

	if err := db.Migrator().CreateConstraint(&models.Category{}, "Category"); err != nil {
		return err
	}

	return db.Migrator().CreateConstraint(&models.Category{}, "fk_categories_parent_categories")
}
