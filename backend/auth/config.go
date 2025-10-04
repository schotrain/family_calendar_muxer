package auth

import (
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

	// JWT secret key - should be set via environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-this-in-production"
	}
	JWTSecret = []byte(secret)

	// Determine if we should use secure connections (HTTPS/SSL)
	useSecure := os.Getenv("USE_SECURE_CONNECTIONS")
	if useSecure == "" {
		useSecure = "false" // Default to false for local development
	}
	UseSecureConnections, _ = strconv.ParseBool(useSecure)
}
