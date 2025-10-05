package main

import (
	"log"
	"net/http"

	"family-calendar-backend/auth"
	"family-calendar-backend/db"
	"family-calendar-backend/rest_api_handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize authentication
	auth.InitAuthConfig()

	// Initialize database
	db.InitDB()

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

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

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
