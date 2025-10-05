package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitAuthConfig_WithValidEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("GOOGLE_REDIRECT_URL", "http://localhost:8080/callback")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("USE_SECURE_CONNECTIONS", "true")
	os.Setenv("ALLOWED_CALLBACKS", "http://localhost:3000/auth/callback,http://localhost:3001/callback")
	defer func() {
		os.Unsetenv("GOOGLE_CLIENT_ID")
		os.Unsetenv("GOOGLE_CLIENT_SECRET")
		os.Unsetenv("GOOGLE_REDIRECT_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("USE_SECURE_CONNECTIONS")
		os.Unsetenv("ALLOWED_CALLBACKS")
	}()

	err := InitAuthConfig()

	assert.NoError(t, err)
	assert.NotNil(t, GoogleOAuthConfig)
	assert.Equal(t, "test-client-id", GoogleOAuthConfig.ClientID)
	assert.Equal(t, "test-client-secret", GoogleOAuthConfig.ClientSecret)
	assert.Equal(t, "http://localhost:8080/callback", GoogleOAuthConfig.RedirectURL)
	assert.Contains(t, GoogleOAuthConfig.Scopes, "https://www.googleapis.com/auth/userinfo.email")
	assert.Contains(t, GoogleOAuthConfig.Scopes, "https://www.googleapis.com/auth/userinfo.profile")

	assert.Equal(t, []byte("test-jwt-secret"), JWTSecret)
	assert.True(t, UseSecureConnections)
	assert.Equal(t, []string{"http://localhost:3000/auth/callback", "http://localhost:3001/callback"}, AllowedCallbacks)
}

func TestInitAuthConfig_DefaultSecureConnections(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Unsetenv("USE_SECURE_CONNECTIONS")
	defer os.Unsetenv("JWT_SECRET")

	err := InitAuthConfig()

	assert.NoError(t, err)
	// Should default to true when not set
	assert.True(t, UseSecureConnections)
}

func TestInitAuthConfig_SecureConnectionsFalse(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("USE_SECURE_CONNECTIONS", "false")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("USE_SECURE_CONNECTIONS")
	}()

	err := InitAuthConfig()

	assert.NoError(t, err)
	assert.False(t, UseSecureConnections)
}

func TestInitAuthConfig_MissingJWTSecret(t *testing.T) {
	// Save and unset JWT_SECRET
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	os.Unsetenv("JWT_SECRET")

	err := InitAuthConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET environment variable is required but not set")
}

func TestInitAuthConfig_InvalidSecureConnectionsValue(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("USE_SECURE_CONNECTIONS", "invalid-value")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("USE_SECURE_CONNECTIONS")
	}()

	err := InitAuthConfig()

	assert.NoError(t, err)
	// Should default to true (secure) when value is invalid
	assert.True(t, UseSecureConnections)
}
