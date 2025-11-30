package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/ItsDobiel/URLShortener/internal/config"
	"github.com/ItsDobiel/URLShortener/internal/shortener"
)

// Handler manages HTTP requests
type Handler struct {
	shortener *shortener.Service
	config    *config.Config
	templates *template.Template
}

// NewHandler creates a new handler instance
func NewHandler(svc *shortener.Service, cfg *config.Config) (*Handler, error) {
	// Parse templates
	tmpl, err := template.ParseGlob(filepath.Join(cfg.TemplatesDir, "*.html"))
	if err != nil {
		return nil, err
	}

	return &Handler{
		shortener: svc,
		config:    cfg,
		templates: tmpl,
	}, nil
}

// HomeHandler displays the main page
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// ShortenHandler processes URL shortening requests
func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		h.renderError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		h.renderError(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Shorten the URL
	shortCode, err := h.shortener.ShortenURL(originalURL)
	if err != nil {
		h.renderError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build short URL
	shortURL := h.config.GetShortURL(shortCode)

	// Render success response
	data := map[string]any{
		"ShortURL":    shortURL,
		"OriginalURL": originalURL,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// RedirectHandler handles short code redirects
func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract short code from path
	shortCode := r.URL.Path[1:] // Remove leading "/"

	if shortCode == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Get original URL
	originalURL, err := h.shortener.GetOriginalURL(shortCode)
	if err != nil {
		h.renderError(w, "Short code not found", http.StatusNotFound)
		return
	}

	// Redirect to original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}

// renderError displays an error page
func (h *Handler) renderError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)

	data := map[string]any{
		"Error":      message,
		"StatusCode": statusCode,
	}

	if err := h.templates.ExecuteTemplate(w, "error.html", data); err != nil {
		http.Error(w, message, statusCode)
	}
}
