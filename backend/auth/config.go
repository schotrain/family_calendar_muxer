package auth

import (
	"log"
	"os"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOAuthConfig *oauth2.Config
	JWTSecret         []byte
	UseSecureConnections bool
)

func InitAuthConfig() {
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"), // e.g., "http://localhost:8080/auth/callback"
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// JWT secret key - MUST be set via environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required but not set")
	}
	JWTSecret = []byte(secret)

	// Determine if we should use secure connections (HTTPS/SSL)
	// Defaults to true for security - explicitly set to false for local development
	useSecure := os.Getenv("USE_SECURE_CONNECTIONS")
	if useSecure == "" {
		useSecure = "true"
	}
	UseSecureConnections, _ = strconv.ParseBool(useSecure)
}
