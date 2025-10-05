package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateFamilyCalendarJWT(t *testing.T) {
	// Set up test secret
	JWTSecret = []byte("test-secret-key")

	tests := []struct {
		name   string
		userID uint
	}{
		{
			name:   "Generate token for user ID 1",
			userID: 1,
		},
		{
			name:   "Generate token for user ID 999",
			userID: 999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateFamilyCalendarJWT(tt.userID)

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// Parse and validate the token
			parsedToken, err := jwt.ParseWithClaims(token, &FamilyCalendarClaims{}, func(token *jwt.Token) (interface{}, error) {
				return JWTSecret, nil
			})

			assert.NoError(t, err)
			assert.True(t, parsedToken.Valid)

			// Check claims
			claims, ok := parsedToken.Claims.(*FamilyCalendarClaims)
			assert.True(t, ok)
			assert.Equal(t, tt.userID, claims.UserID)
			assert.Equal(t, "family-calendar-backend", claims.Issuer)
			assert.NotNil(t, claims.ExpiresAt)
			assert.NotNil(t, claims.IssuedAt)

			// Check expiration is approximately 24 hours from now
			expectedExpiry := time.Now().Add(24 * time.Hour)
			assert.WithinDuration(t, expectedExpiry, claims.ExpiresAt.Time, 2*time.Second)
		})
	}
}

func TestFamilyCalendarClaims_Expiration(t *testing.T) {
	JWTSecret = []byte("test-secret-key")

	// Create an expired token
	expiredClaims := FamilyCalendarClaims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-25 * time.Hour)),
			Issuer:    "family-calendar-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenString, err := token.SignedString(JWTSecret)
	assert.NoError(t, err)

	// Try to parse expired token
	parsedToken, err := jwt.ParseWithClaims(tokenString, &FamilyCalendarClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	assert.Error(t, err)
	assert.False(t, parsedToken.Valid)
}

func TestFamilyCalendarClaims_InvalidSignature(t *testing.T) {
	JWTSecret = []byte("test-secret-key")

	// Create token with one secret
	token, err := GenerateFamilyCalendarJWT(1)
	assert.NoError(t, err)

	// Try to parse with different secret
	differentSecret := []byte("different-secret")
	parsedToken, err := jwt.ParseWithClaims(token, &FamilyCalendarClaims{}, func(token *jwt.Token) (interface{}, error) {
		return differentSecret, nil
	})

	assert.Error(t, err)
	assert.NotNil(t, parsedToken)
	assert.False(t, parsedToken.Valid)
}
