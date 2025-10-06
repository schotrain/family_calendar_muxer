package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorsMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		allowedOrigin  string
		method         string
		expectedOrigin string
		expectedStatus int
	}{
		{
			name:           "GET request with default origin",
			allowedOrigin:  "http://localhost:3000",
			method:         "GET",
			expectedOrigin: "http://localhost:3000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "OPTIONS preflight request",
			allowedOrigin:  "http://localhost:3000",
			method:         "OPTIONS",
			expectedOrigin: "http://localhost:3000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST request with custom origin",
			allowedOrigin:  "https://example.com",
			method:         "POST",
			expectedOrigin: "https://example.com",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with CORS middleware
			wrappedHandler := corsMiddleware(tt.allowedOrigin)(handler)

			// Create request
			req := httptest.NewRequest(tt.method, "/test", nil)
			rr := httptest.NewRecorder()

			// Serve
			wrappedHandler.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Assert CORS headers
			assert.Equal(t, tt.expectedOrigin, rr.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
			assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
		})
	}
}

func TestCorsMiddleware_OptionsDoesNotCallNext(t *testing.T) {
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := corsMiddleware("http://localhost:3000")(handler)

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	// Handler should not be called for OPTIONS requests
	assert.False(t, handlerCalled)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCorsMiddleware_NonOptionsCallsNext(t *testing.T) {
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := corsMiddleware("http://localhost:3000")(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	// Handler should be called for non-OPTIONS requests
	assert.True(t, handlerCalled)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSetupRouter_Success(t *testing.T) {
	// Set required environment variables
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	t.Setenv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback")
	t.Setenv("DB_TYPE", "sqlite")

	router, err := setupRouter()

	assert.NoError(t, err)
	assert.NotNil(t, router)

	// Clean up database
	// Note: In a real test you might want to close the database connection
}

func TestSetupRouter_WithCorsOrigin(t *testing.T) {
	// Set required environment variables
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	t.Setenv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback")
	t.Setenv("DB_TYPE", "sqlite")
	t.Setenv("CORS_ALLOWED_ORIGIN", "https://example.com")

	router, err := setupRouter()

	assert.NoError(t, err)
	assert.NotNil(t, router)

	// Test that CORS origin is set correctly by making a request
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, "https://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestSetupRouter_AuthInitError(t *testing.T) {
	// Don't set JWT_SECRET to trigger auth init error
	t.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	t.Setenv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback")

	router, err := setupRouter()

	assert.Error(t, err)
	assert.Nil(t, router)
	assert.Contains(t, err.Error(), "JWT_SECRET")
}

func TestSetupRouter_DBInitError(t *testing.T) {
	// Set auth env vars but use invalid DB type
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	t.Setenv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback")
	t.Setenv("DB_TYPE", "invalid-db-type")

	router, err := setupRouter()

	assert.Error(t, err)
	assert.Nil(t, router)
	assert.Contains(t, err.Error(), "unsupported database type")
}

