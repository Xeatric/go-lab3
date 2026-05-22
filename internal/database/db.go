package database

import (
	"fmt"
	"log"

	"paving-tiles-api/internal/config"
	"paving-tiles-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// RunMigrations - автоматическое создание таблиц
func RunMigrations(db *gorm.DB) error {
	log.Println("Running migrations...")

	// Миграция для таблицы users
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate users: %v", err)
	}
	log.Println("Users table migrated")

	// Миграция для таблицы tokens
	if err := db.AutoMigrate(&models.Token{}); err != nil {
		return fmt.Errorf("failed to migrate tokens: %v", err)
	}
	log.Println("Tokens table migrated")

	// Миграция для таблицы tiles (существующая)
	if err := db.AutoMigrate(&models.Tile{}); err != nil {
		return fmt.Errorf("failed to migrate tiles: %v", err)
	}
	log.Println("Tiles table migrated")

	log.Println("All migrations completed successfully")
	return nil
}
