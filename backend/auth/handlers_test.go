package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"family-calendar-backend/db/models"
	"family-calendar-backend/db/services"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func setupAuthTests() {
	// Set required environment variables for tests
	os.Setenv("JWT_SECRET", "test-secret-key-for-auth-tests")
	os.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback")
	os.Setenv("USE_SECURE_CONNECTIONS", "false")
	os.Setenv("ALLOWED_CALLBACKS", "http://localhost:3000/auth/callback")
	InitAuthConfig()
}

func TestLoginHandler(t *testing.T) {
	setupAuthTests()

	req := httptest.NewRequest("GET", "/auth/google", nil)
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	// Should redirect to Google OAuth
	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
	assert.NotEmpty(t, rr.Header().Get("Location"))
	assert.Contains(t, rr.Header().Get("Location"), "accounts.google.com/o/oauth2")

	// Should set state cookie
	cookies := rr.Result().Cookies()
	var stateCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "oauth_state" {
			stateCookie = cookie
			break
		}
	}
	assert.NotNil(t, stateCookie)
	assert.Equal(t, "oauth_state", stateCookie.Name)
	assert.NotEmpty(t, stateCookie.Value)
	assert.Equal(t, 300, stateCookie.MaxAge)
	assert.True(t, stateCookie.HttpOnly)
	assert.Equal(t, http.SameSiteLaxMode, stateCookie.SameSite)
}

func TestLoginHandler_WithValidCallback(t *testing.T) {
	setupAuthTests()

	req := httptest.NewRequest("GET", "/auth/google?callback=http://localhost:3000/auth/callback", nil)
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	// Should redirect to Google OAuth
	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
	assert.Contains(t, rr.Header().Get("Location"), "accounts.google.com/o/oauth2")

	// Should set both state and callback cookies
	cookies := rr.Result().Cookies()
	var stateCookie, callbackCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "oauth_state" {
			stateCookie = cookie
		}
		if cookie.Name == "oauth_callback" {
			callbackCookie = cookie
		}
	}
	assert.NotNil(t, stateCookie)
	assert.NotNil(t, callbackCookie)
	assert.Equal(t, "http://localhost:3000/auth/callback", callbackCookie.Value)
	assert.Equal(t, 300, callbackCookie.MaxAge)
	assert.True(t, callbackCookie.HttpOnly)
}

func TestLoginHandler_WithInvalidCallback(t *testing.T) {
	setupAuthTests()

	req := httptest.NewRequest("GET", "/auth/google?callback=http://evil.com/steal", nil)
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Contains(t, rr.Body.String(), "callback URL is not allowed")
}

func TestCallbackHandler_MissingStateCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "/auth/google/callback?state=test&code=test", nil)
	rr := httptest.NewRecorder()

	CallbackHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "State cookie not found")
}

func TestCallbackHandler_InvalidState(t *testing.T) {
	req := httptest.NewRequest("GET", "/auth/google/callback?state=wrong&code=test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "correct",
	})
	rr := httptest.NewRecorder()

	CallbackHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid state parameter")
}

func TestGenerateStateToken(t *testing.T) {
	token1 := generateStateToken()
	token2 := generateStateToken()

	// Tokens should be non-empty
	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)

	// Tokens should be different (random)
	assert.NotEqual(t, token1, token2)

	// Token should be base64 encoded
	assert.True(t, len(token1) > 40) // 32 bytes encoded should be > 40 chars
}

func TestRenderTokenPage_TemplateNotFound(t *testing.T) {
	rr := httptest.NewRecorder()
	userInfo := GoogleUserInfo{
		GivenName:  "John",
		FamilyName: "Doe",
		Email:      "john@example.com",
	}

	// This will fail because the template path won't exist in test environment
	// But we can test it handles the error gracefully
	renderTokenPage(rr, "test-token", userInfo)

	// Should return 500 if template fails to load
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRenderTokenPage_Success(t *testing.T) {
	// Create a temporary template file
	tmpDir := t.TempDir()
	authDir := filepath.Join(tmpDir, "auth", "templates")
	err := os.MkdirAll(authDir, 0755)
	assert.NoError(t, err)

	templateContent := `<!DOCTYPE html>
<html>
<head><title>Auth Success</title></head>
<body>
<h1>Hello {{.GivenName}} {{.FamilyName}}</h1>
<p>Email: {{.Email}}</p>
<p>Token: {{.Token}}</p>
</body>
</html>`

	templateFile := filepath.Join(authDir, "auth_success.html")
	err = os.WriteFile(templateFile, []byte(templateContent), 0644)
	assert.NoError(t, err)

	// Change to temp directory for template resolution
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Test successful rendering
	rr := httptest.NewRecorder()
	userInfo := GoogleUserInfo{
		GivenName:  "John",
		FamilyName: "Doe",
		Email:      "john@example.com",
	}

	renderTokenPage(rr, "test-jwt-token", userInfo)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"))
	assert.Contains(t, rr.Body.String(), "Hello John Doe")
	assert.Contains(t, rr.Body.String(), "john@example.com")
	assert.Contains(t, rr.Body.String(), "test-jwt-token")
}

func TestRenderTokenPage_TemplateExecutionError(t *testing.T) {
	// Create a template that will fail during execution
	tmpDir := t.TempDir()
	authDir := filepath.Join(tmpDir, "auth", "templates")
	err := os.MkdirAll(authDir, 0755)
	assert.NoError(t, err)

	// Template with undefined field that will cause execution error
	templateContent := `{{.UndefinedField.Method}}`
	templateFile := filepath.Join(authDir, "auth_success.html")
	err = os.WriteFile(templateFile, []byte(templateContent), 0644)
	assert.NoError(t, err)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	rr := httptest.NewRecorder()
	userInfo := GoogleUserInfo{
		GivenName:  "Test",
		FamilyName: "User",
		Email:      "test@example.com",
	}

	renderTokenPage(rr, "token", userInfo)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to render template")
}

func TestCallbackHandler_MissingCode(t *testing.T) {
	req := httptest.NewRequest("GET", "/auth/google/callback?state=test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test",
	})
	rr := httptest.NewRecorder()

	CallbackHandler(rr, req)

	// Without code parameter, the OAuth exchange will fail
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to exchange token")
}

func TestGoogleUserInfo_Structure(t *testing.T) {
	// Test that GoogleUserInfo struct has expected fields
	userInfo := GoogleUserInfo{
		ID:            "123456",
		Sub:           "123456",
		Email:         "test@example.com",
		Name:          "Test User",
		GivenName:     "Test",
		FamilyName:    "User",
		Picture:       "https://example.com/pic.jpg",
		VerifiedEmail: true,
	}

	assert.Equal(t, "123456", userInfo.ID)
	assert.Equal(t, "123456", userInfo.Sub)
	assert.Equal(t, "test@example.com", userInfo.Email)
	assert.Equal(t, "Test User", userInfo.Name)
	assert.Equal(t, "Test", userInfo.GivenName)
	assert.Equal(t, "User", userInfo.FamilyName)
	assert.Equal(t, "https://example.com/pic.jpg", userInfo.Picture)
	assert.True(t, userInfo.VerifiedEmail)
}

func TestGoogleUserInfo_GetUserID(t *testing.T) {
	tests := []struct {
		name     string
		userInfo GoogleUserInfo
		expected string
	}{
		{
			name: "ID field present (v2 API)",
			userInfo: GoogleUserInfo{
				ID:  "google-id-123",
				Sub: "",
			},
			expected: "google-id-123",
		},
		{
			name: "Sub field present (v3 API)",
			userInfo: GoogleUserInfo{
				ID:  "",
				Sub: "google-sub-456",
			},
			expected: "google-sub-456",
		},
		{
			name: "Both fields present (ID takes precedence)",
			userInfo: GoogleUserInfo{
				ID:  "google-id-123",
				Sub: "google-sub-456",
			},
			expected: "google-id-123",
		},
		{
			name: "Neither field present",
			userInfo: GoogleUserInfo{
				ID:  "",
				Sub: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userInfo.GetUserID()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoginHandler_StateInURL(t *testing.T) {
	setupAuthTests()

	req := httptest.NewRequest("GET", "/auth/google", nil)
	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	// Get the state from cookie
	cookies := rr.Result().Cookies()
	var stateCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "oauth_state" {
			stateCookie = cookie
			break
		}
	}

	// Verify state is in the redirect URL (it will be URL encoded)
	location := rr.Header().Get("Location")
	assert.Contains(t, location, "state=")
	assert.NotEmpty(t, stateCookie.Value)
}

func TestCallbackHandler_StateValidation(t *testing.T) {
	tests := []struct {
		name           string
		stateParam     string
		cookieValue    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Matching state",
			stateParam:     "abc123",
			cookieValue:    "abc123",
			expectedStatus: http.StatusInternalServerError, // Will fail on exchange, but passes state check
			expectedBody:   "Failed to exchange token",
		},
		{
			name:           "Mismatched state",
			stateParam:     "abc123",
			cookieValue:    "xyz789",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid state parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/auth/google/callback?state="+tt.stateParam+"&code=testcode", nil)
			req.AddCookie(&http.Cookie{
				Name:  "oauth_state",
				Value: tt.cookieValue,
			})
			rr := httptest.NewRecorder()

			CallbackHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, strings.TrimSpace(rr.Body.String()), tt.expectedBody)
		})
	}
}

// TestGoogleUserInfoJSONParsing tests JSON parsing of Google user info
func TestGoogleUserInfoJSONParsing(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expectError bool
	}{
		{
			name: "Valid JSON with ID (v2 API)",
			jsonData: `{
				"id": "123456",
				"email": "test@example.com",
				"name": "Test User",
				"given_name": "Test",
				"family_name": "User",
				"picture": "https://example.com/pic.jpg",
				"verified_email": true
			}`,
			expectError: false,
		},
		{
			name: "Valid JSON with Sub (v3 API)",
			jsonData: `{
				"sub": "123456",
				"email": "test@example.com",
				"name": "Test User",
				"given_name": "Test",
				"family_name": "User",
				"picture": "https://example.com/pic.jpg",
				"verified_email": true
			}`,
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			jsonData:    `{"id": "123", "email": }`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var userInfo GoogleUserInfo
			err := json.NewDecoder(bytes.NewBufferString(tt.jsonData)).Decode(&userInfo)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test@example.com", userInfo.Email)
				assert.Equal(t, "Test", userInfo.GivenName)
			}
		})
	}
}

// TestTemplateDataStructure tests the template data structure
func TestTemplateDataStructure(t *testing.T) {
	// Test that the template data structure can be marshaled/unmarshaled
	data := struct {
		Token      string
		GivenName  string
		FamilyName string
		Email      string
	}{
		Token:      "test-token",
		GivenName:  "John",
		FamilyName: "Doe",
		Email:      "john@example.com",
	}

	// Test with actual template parsing
	tmpl, err := template.New("test").Parse("{{.Token}} {{.GivenName}} {{.FamilyName}} {{.Email}}")
	assert.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	assert.NoError(t, err)
	assert.Equal(t, "test-token John Doe john@example.com", buf.String())
}

// TestResponseBodyClosure tests response body handling
func TestResponseBodyClosure(t *testing.T) {
	// Create a ReadCloser
	rc := io.NopCloser(bytes.NewBufferString("test"))

	// Simulate reading and closing
	data, err := io.ReadAll(rc)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(data))

	err = rc.Close()
	assert.NoError(t, err)
}

func TestCallbackHandler_SuccessfulOAuthFlow(t *testing.T) {
	// Save original functions
	originalExchange := exchangeToken
	originalGetUserInfo := getUserInfo
	originalFindOrCreate := services.FindOrCreateUser
	defer func() {
		exchangeToken = originalExchange
		getUserInfo = originalGetUserInfo
		services.FindOrCreateUser = originalFindOrCreate
	}()

	// Mock successful OAuth token exchange
	exchangeToken = func(ctx context.Context, code string) (*oauth2.Token, error) {
		assert.Equal(t, "test-code", code)
		return &oauth2.Token{AccessToken: "mock-access-token"}, nil
	}

	// Mock successful user info retrieval
	getUserInfo = func(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
		return &GoogleUserInfo{
			ID:         "google-user-123",
			Email:      "test@example.com",
			GivenName:  "Test",
			FamilyName: "User",
		}, nil
	}

	// Mock database user creation
	services.FindOrCreateUser = func(provider, providerID, givenName, familyName, email string) (*models.User, error) {
		return &models.User{
			GivenName:      givenName,
			FamilyName:     familyName,
			Email:          email,
			AuthProvider:   provider,
			AuthProviderID: providerID,
		}, nil
	}

	// Create temporary template
	tmpDir := t.TempDir()
	authDir := filepath.Join(tmpDir, "auth", "templates")
	err := os.MkdirAll(authDir, 0755)
	assert.NoError(t, err)

	templateContent := `Token: {{.Token}}`
	templateFile := filepath.Join(authDir, "auth_success.html")
	err = os.WriteFile(templateFile, []byte(templateContent), 0644)
	assert.NoError(t, err)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Setup request
	req := httptest.NewRequest("GET", "/auth/google/callback?state=test-state&code=test-code", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test-state",
	})
	rr := httptest.NewRecorder()

	CallbackHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Token:")
}

func TestCallbackHandler_ExchangeTokenError(t *testing.T) {
	// Save original function
	originalExchange := exchangeToken
	defer func() { exchangeToken = originalExchange }()

	// Mock failed OAuth token exchange
	exchangeToken = func(ctx context.Context, code string) (*oauth2.Token, error) {
		return nil, assert.AnError
	}

	req := httptest.NewRequest("GET", "/auth/google/callback?state=test&code=test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test",
	})
	rr := httptest.NewRecorder()

	CallbackHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to exchange token")
}

func TestCallbackHandler_GetUserInfoError(t *testing.T) {
	// Save original functions
	originalExchange := exchangeToken
	originalGetUserInfo := getUserInfo
	defer func() {
		exchangeToken = originalExchange
		getUserInfo = originalGetUserInfo
	}()

	// Mock successful exchange but failed user info
	exchangeToken = func(ctx context.Context, code string) (*oauth2.Token, error) {
		return &oauth2.Token{AccessToken: "mock-token"}, nil
	}

	getUserInfo = func(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
		return nil, assert.AnError
	}

	req := httptest.NewRequest("GET", "/auth/google/callback?state=test&code=test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test",
	})
	rr := httptest.NewRecorder()

	CallbackHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to get user info")
}

func TestCallbackHandler_DatabaseError(t *testing.T) {
	// Save original functions
	originalExchange := exchangeToken
	originalGetUserInfo := getUserInfo
	originalFindOrCreate := services.FindOrCreateUser
	defer func() {
		exchangeToken = originalExchange
		getUserInfo = originalGetUserInfo
		services.FindOrCreateUser = originalFindOrCreate
	}()

	// Mock successful OAuth flow
	exchangeToken = func(ctx context.Context, code string) (*oauth2.Token, error) {
		return &oauth2.Token{AccessToken: "mock-token"}, nil
	}

	getUserInfo = func(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
		return &GoogleUserInfo{
			ID:         "google-123",
			Email:      "test@example.com",
			GivenName:  "Test",
			FamilyName: "User",
		}, nil
	}

	// Mock database error
	services.FindOrCreateUser = func(provider, providerID, givenName, familyName, email string) (*models.User, error) {
		return nil, assert.AnError
	}

	req := httptest.NewRequest("GET", "/auth/google/callback?state=test&code=test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test",
	})
	rr := httptest.NewRecorder()

	CallbackHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to process user")
}

func TestCallbackHandler_WithCallbackRedirect(t *testing.T) {
	setupAuthTests()

	// Save original functions
	originalExchange := exchangeToken
	originalGetUserInfo := getUserInfo
	originalFindOrCreate := services.FindOrCreateUser
	defer func() {
		exchangeToken = originalExchange
		getUserInfo = originalGetUserInfo
		services.FindOrCreateUser = originalFindOrCreate
	}()

	// Mock successful OAuth flow
	exchangeToken = func(ctx context.Context, code string) (*oauth2.Token, error) {
		return &oauth2.Token{AccessToken: "mock-access-token"}, nil
	}

	getUserInfo = func(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
		return &GoogleUserInfo{
			ID:         "google-user-123",
			Email:      "test@example.com",
			GivenName:  "Test",
			FamilyName: "User",
		}, nil
	}

	services.FindOrCreateUser = func(provider, providerID, givenName, familyName, email string) (*models.User, error) {
		user := &models.User{
			GivenName:      givenName,
			FamilyName:     familyName,
			Email:          email,
			AuthProvider:   provider,
			AuthProviderID: providerID,
		}
		user.ID = 1
		return user, nil
	}

	// Setup request with callback cookie
	req := httptest.NewRequest("GET", "/auth/google/callback?state=test-state&code=test-code", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test-state",
	})
	req.AddCookie(&http.Cookie{
		Name:  "oauth_callback",
		Value: "http://localhost:3000/auth/callback",
	})
	rr := httptest.NewRecorder()

	CallbackHandler(rr, req)

	// Should redirect to callback URL with token
	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
	location := rr.Header().Get("Location")
	assert.Contains(t, location, "http://localhost:3000/auth/callback?token=")
	assert.NotContains(t, location, "token=http") // Token shouldn't contain URL

	// Should clear cookies
	cookies := rr.Result().Cookies()
	var stateCookie, callbackCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "oauth_state" {
			stateCookie = cookie
		}
		if cookie.Name == "oauth_callback" {
			callbackCookie = cookie
		}
	}
	assert.NotNil(t, stateCookie)
	assert.Equal(t, -1, stateCookie.MaxAge)
	assert.NotNil(t, callbackCookie)
	assert.Equal(t, -1, callbackCookie.MaxAge)
}

func TestCallbackHandler_MissingUserID(t *testing.T) {
	setupAuthTests()

	// Save original functions
	originalExchange := exchangeToken
	originalGetUserInfo := getUserInfo
	defer func() {
		exchangeToken = originalExchange
		getUserInfo = originalGetUserInfo
	}()

	// Mock OAuth flow with empty user ID
	exchangeToken = func(ctx context.Context, code string) (*oauth2.Token, error) {
		return &oauth2.Token{AccessToken: "mock-access-token"}, nil
	}

	getUserInfo = func(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
		return &GoogleUserInfo{
			ID:         "", // Empty ID
			Sub:        "", // Empty Sub
			Email:      "test@example.com",
			GivenName:  "Test",
			FamilyName: "User",
		}, nil
	}

	req := httptest.NewRequest("GET", "/auth/google/callback?state=test-state&code=test-code", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test-state",
	})
	rr := httptest.NewRecorder()

	CallbackHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid user info from Google")
}

