package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort      string
	ServerHost      string
	ShortDomain     string
	DatabasePath    string
	ShortCodeLength int
	TemplatesDir    string
}

// Load reads configuration from environment variables
// It returns an error if required variables are missing or invalid
func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		ServerHost:   getEnv("SERVER_HOST", "localhost"),
		ShortDomain:  getEnv("SHORT_DOMAIN", "localhost:8080"),
		DatabasePath: getEnv("DATABASE_PATH", "./database"),
		TemplatesDir: getEnv("TEMPLATES_DIR", "templates"),
	}

	lengthStr := getEnv("SHORT_CODE_LENGTH", "7")
	length, err := strconv.Atoi(lengthStr)
	if err != nil || length < 4 || length > 12 {
		return nil, fmt.Errorf("invalid SHORT_CODE_LENGTH: must be between 4 and 12")
	}
	config.ShortCodeLength = length

	return config, nil
}

// GetAddress returns the full server address for binding
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)
}

// GetShortURL constructs the full short URL from a short code
func (c *Config) GetShortURL(shortCode string) string {
	return fmt.Sprintf("http://%s/%s", c.ShortDomain, shortCode)
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
