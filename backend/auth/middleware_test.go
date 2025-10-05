package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestRequireAuth_ValidToken(t *testing.T) {
	JWTSecret = []byte("test-secret-key")

	// Create a valid token
	token, err := GenerateFamilyCalendarJWT(123)
	assert.NoError(t, err)

	// Create test handler
	var capturedUserID uint
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())
		assert.True(t, ok)
		capturedUserID = userID
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with auth middleware
	handler := RequireAuth(testHandler)

	// Create request with valid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, uint(123), capturedUserID)
}

func TestRequireAuth_MissingAuthHeader(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handler should not be called")
	})

	handler := RequireAuth(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Authorization header required")
}

func TestRequireAuth_InvalidAuthHeaderFormat(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{
			name:   "Missing Bearer prefix",
			header: "token-without-bearer",
		},
		{
			name:   "Wrong prefix",
			header: "Basic some-token",
		},
		{
			name:   "Too many parts",
			header: "Bearer token extra",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("Handler should not be called")
			})

			handler := RequireAuth(testHandler)

			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.header)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		})
	}
}

func TestRequireAuth_ExpiredToken(t *testing.T) {
	JWTSecret = []byte("test-secret-key")

	// Create expired token
	expiredClaims := FamilyCalendarClaims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-25 * time.Hour)),
			Issuer:    "family-calendar-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenString, err := token.SignedString(JWTSecret)
	assert.NoError(t, err)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handler should not be called")
	})

	handler := RequireAuth(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid or expired token")
}

func TestRequireAuth_InvalidSignature(t *testing.T) {
	JWTSecret = []byte("test-secret-key")

	// Create token with different secret
	differentSecret := []byte("different-secret")
	claims := FamilyCalendarClaims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "family-calendar-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(differentSecret)
	assert.NoError(t, err)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handler should not be called")
	})

	handler := RequireAuth(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestRequireAuth_WrongSigningMethod(t *testing.T) {
	JWTSecret = []byte("test-secret-key")

	// Create token with RS256 instead of HS256
	claims := FamilyCalendarClaims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "family-calendar-backend",
		},
	}

	// Use none algorithm (not recommended but for testing)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handler should not be called")
	})

	handler := RequireAuth(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name          string
		contextValue  interface{}
		expectedID    uint
		expectedFound bool
	}{
		{
			name:          "Valid user ID in context",
			contextValue:  uint(42),
			expectedID:    42,
			expectedFound: true,
		},
		{
			name:          "No user ID in context",
			contextValue:  nil,
			expectedID:    0,
			expectedFound: false,
		},
		{
			name:          "Wrong type in context",
			contextValue:  "not-a-uint",
			expectedID:    0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.contextValue != nil {
				ctx = context.WithValue(ctx, UserIDContextKey, tt.contextValue)
			}

			userID, found := GetUserIDFromContext(ctx)

			assert.Equal(t, tt.expectedFound, found)
			assert.Equal(t, tt.expectedID, userID)
		})
	}
}

func TestSetUserIDInContext(t *testing.T) {
	// Test SetUserIDInContext for coverage
	ctx := context.Background()
	newCtx := SetUserIDInContext(ctx, 123)
	userID, ok := GetUserIDFromContext(newCtx)
	assert.True(t, ok)
	assert.Equal(t, uint(123), userID)
}
