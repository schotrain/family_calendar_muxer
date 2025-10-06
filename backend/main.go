package main

import (
	"log"
	"net/http"
	"os"

	"family-calendar-backend/auth"
	"family-calendar-backend/db"
	"family-calendar-backend/rest_api_handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func corsMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func setupRouter() (*chi.Mux, error) {
	// Initialize authentication
	if err := auth.InitAuthConfig(); err != nil {
		return nil, err
	}

	// Initialize database
	if err := db.InitDB(); err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	// Get CORS allowed origin from environment
	corsOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:3000" // Default for development
	}

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(corsMiddleware(corsOrigin))

	// Auth routes (not part of REST API)
	r.Get("/auth/google", auth.LoginHandler)
	r.Get("/auth/google/callback", auth.CallbackHandler)

	// Public REST API routes (no authentication required)
	r.Get("/health", rest_api_handlers.HealthCheck)

	// Protected REST API routes (authentication required)
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth)
		r.Get("/api/userinfo", rest_api_handlers.UserInfo)
	})

	return r, nil
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	r, err := setupRouter()
	if err != nil {
		log.Printf("Failed to setup router: %v", err)
		return
	}

	log.Println("Server starting on 0.0.0.0:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Printf("Server stopped: %v", err)
	}
}
