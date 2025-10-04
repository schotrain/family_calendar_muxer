package main

import (
	"log"
	"net/http"

	"family-calendar-backend/auth"
	"family-calendar-backend/database"
	"family-calendar-backend/rest_api_handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize authentication
	auth.InitAuthConfig()

	// Initialize database
	database.InitDB()

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Auth routes (not part of REST API)
	r.Get("/auth", auth.LoginHandler)
	r.Get("/auth/callback", auth.CallbackHandler)

	// REST API routes
	r.Get("/health", rest_api_handlers.HealthCheck)
	r.Post("/api/users", rest_api_handlers.CreateUser)
	r.Get("/api/users/{id}", rest_api_handlers.GetUser)

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
