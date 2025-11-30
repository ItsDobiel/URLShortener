package router

import (
	"net/http"

	"github.com/ItsDobiel/URLShortener/internal/handlers"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(handler *handlers.Handler, templatesDir string) *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files (In this case CSS files)
	fs := http.FileServer(http.Dir(templatesDir + "/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Home page - displays the URL shortening form
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If path is exactly "/", show home page
		if r.URL.Path == "/" {
			handler.HomeHandler(w, r)
			return
		}

		// Otherwise, treat as short code redirect
		handler.RedirectHandler(w, r)
	})

	// Shorten endpoint - processes URL shortening requests
	mux.HandleFunc("/shorten", handler.ShortenHandler)

	return mux
}
