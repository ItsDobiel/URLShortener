package shortener

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/ItsDobiel/URLShortener/internal/database"
	"github.com/ItsDobiel/URLShortener/internal/models"
)

const (
	// maxCollisionRetries defines how many times to retry on hash collision
	maxCollisionRetries = 5
)

// Service handles URL shortening operations
type Service struct {
	codeLength int
}

// NewService creates a new shortener service
func NewService(codeLength int) *Service {
	return &Service{
		codeLength: codeLength,
	}
}

// ShortenURL creates a short code for the given URL
// If the URL has been shortened before, it returns the existing short code
// Returns the short code and any error encountered
func (s *Service) ShortenURL(rawURL string) (string, error) {
	// Validate URL format
	if err := s.validateURL(rawURL); err != nil {
		return "", err
	}

	// Normalize the URL for consistent handling
	normalizedURL := s.normalizeURL(rawURL)

	existingURL, err := database.FindByNormalizedURL(normalizedURL)
	if err == nil && existingURL != nil {
		return existingURL.ShortCode, nil
	}

	shortCode, err := s.generateUniqueShortCode(normalizedURL)
	if err != nil {
		return "", err
	}

	urlModel := &models.URL{
		ShortCode:     shortCode,
		OriginalURL:   rawURL,
		NormalizedURL: normalizedURL,
	}

	if err := database.Create(urlModel); err != nil {
		return "", fmt.Errorf("failed to save URL: %w", err)
	}

	return shortCode, nil
}

// GetOriginalURL retrieves the original URL for a given short code
func (s *Service) GetOriginalURL(shortCode string) (string, error) {
	if !s.isValidShortCode(shortCode) {
		return "", fmt.Errorf("invalid short code format")
	}

	urlModel, err := database.FindByShortCode(shortCode)
	if err != nil {
		return "", fmt.Errorf("short code not found")
	}

	return urlModel.OriginalURL, nil
}

// validateURL checks if the URL is valid and uses supported protocol
func (s *Service) validateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}

	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("only HTTP and HTTPS protocols are supported")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a valid host")
	}

	return nil
}

// normalizeURL converts a URL to a canonical form for duplicate detection
// This handles trailing slashes, case-insensitive protocols/domains.
func (s *Service) normalizeURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)

	parsedURL.Host = strings.ToLower(parsedURL.Host)

	if len(parsedURL.Path) > 1 && strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")
	}

	if parsedURL.Path == "" {
		parsedURL.Path = "/"
	}

	return parsedURL.String()
}

// generateUniqueShortCode creates a short code that doesn't collide with existing ones
func (s *Service) generateUniqueShortCode(normalizedURL string) (string, error) {
	for attempt := 0; attempt < maxCollisionRetries; attempt++ {
		shortCode := s.generateShortCode(normalizedURL, attempt)

		taken, err := database.IsShortCodeTaken(shortCode)
		if err != nil {
			return "", fmt.Errorf("failed to check short code availability: %w", err)
		}

		if !taken {
			return shortCode, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique short code after %d attempts", maxCollisionRetries)
}

// generateShortCode creates a short code using SHA-256 hash
// The attempt parameter allows generating different codes on collision
func (s *Service) generateShortCode(normalizedURL string, attempt int) string {
	// Add attempt number to input for different hashes on collision
	input := fmt.Sprintf("%s:%d", normalizedURL, attempt)

	hash := sha256.Sum256([]byte(input))

	encoded := base64.URLEncoding.EncodeToString(hash[:])

	encoded = strings.TrimRight(encoded, "=")

	if len(encoded) > s.codeLength {
		return encoded[:s.codeLength]
	}

	return encoded
}

// isValidShortCode checks if a short code matches expected format
func (s *Service) isValidShortCode(shortCode string) bool {
	if len(shortCode) < 4 || len(shortCode) > 20 {
		return false
	}

	// Checks that it only contains base64 URL-safe characters
	for _, char := range shortCode {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}

	return true
}
