package database

import (
	"fmt"

	"github.com/ItsDobiel/URLShortener/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize sets up the database connection and performs migrations
// It returns an error if the connection fails or migrations fail
func Initialize(dbPath string) error {
	var err error

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := DB.AutoMigrate(&models.URL{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// FindByShortCode retrieves a URL by its short code
func FindByShortCode(shortCode string) (*models.URL, error) {
	var url models.URL
	result := DB.Where("short_code = ?", shortCode).First(&url)
	if result.Error != nil {
		return nil, result.Error
	}
	return &url, nil
}

// FindByNormalizedURL retrieves a URL by its normalized form
// This is used to check for duplicate URLs
func FindByNormalizedURL(normalizedURL string) (*models.URL, error) {
	var url models.URL
	result := DB.Where("normalized_url = ?", normalizedURL).First(&url)
	if result.Error != nil {
		return nil, result.Error
	}
	return &url, nil
}

// Create saves a new URL mapping to the database
func Create(url *models.URL) error {
	result := DB.Create(url)
	return result.Error
}

// IsShortCodeTaken checks if a short code already exists
func IsShortCodeTaken(shortCode string) (bool, error) {
	var count int64
	result := DB.Model(&models.URL{}).Where("short_code = ?", shortCode).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
