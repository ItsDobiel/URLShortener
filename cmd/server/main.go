package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ItsDobiel/URLShortener/internal/config"
	"github.com/ItsDobiel/URLShortener/internal/database"
	"github.com/ItsDobiel/URLShortener/internal/handlers"
	"github.com/ItsDobiel/URLShortener/internal/router"
	"github.com/ItsDobiel/URLShortener/internal/shortener"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if _, err := os.Stat(cfg.DatabasePath); os.IsNotExist(err) {
		// Directory doesn't exist, create it
		err := os.MkdirAll(cfg.DatabasePath, 0755)
		if err != nil {
			log.Fatalf("failed to create directory: %v", err)
		}
		log.Printf("Directory created: %s\n", cfg.DatabasePath)
	}

	if err := database.Initialize(cfg.DatabasePath + "/urlshortener.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	log.Println("Database initialized successfully")

	svc := shortener.NewService(cfg.ShortCodeLength)

	handler, err := handlers.NewHandler(svc, cfg)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	mux := router.SetupRouter(handler, cfg.TemplatesDir)

	server := &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: mux,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		if err := server.Close(); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
	}()

	log.Printf("Server starting on http://%s", cfg.GetAddress())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("Server stopped")
}
