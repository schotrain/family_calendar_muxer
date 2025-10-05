package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"family-calendar-backend/db/services"

	"golang.org/x/oauth2"
)

// GoogleUserInfo represents the user information from Google
type GoogleUserInfo struct {
	Sub           string `json:"sub"` // Unique Google user ID
	Email         string `json:"email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// OAuth function variables for testing.
// NOTE: The default implementations below make actual network calls to Google's OAuth API
// and are intentionally NOT covered by unit tests. They are mocked in all tests to avoid
// network dependencies. These functions are only used in production.
var (
	exchangeToken = func(ctx context.Context, code string) (*oauth2.Token, error) {
		return GoogleOAuthConfig.Exchange(ctx, code)
	}
	getUserInfo = func(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
		client := GoogleOAuthConfig.Client(ctx, token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var userInfo GoogleUserInfo
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			return nil, err
		}
		return &userInfo, nil
	}
)

// LoginHandler initiates the OAuth flow
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	state := generateStateToken()

	// Store state in session/cookie (simplified for now)
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		Secure:   UseSecureConnections,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to Google's OAuth consent page
	url := GoogleOAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// CallbackHandler handles the OAuth callback from Google
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Verify state token
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "State cookie not found", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("state") != stateCookie.Value {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Exchange authorization code for token
	code := r.URL.Query().Get("code")
	token, err := exchangeToken(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get user info from Google
	userInfo, err := getUserInfo(context.Background(), token)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Find or create user in database
	user, err := services.FindOrCreateUser("google", userInfo.Sub, userInfo.GivenName, userInfo.FamilyName, userInfo.Email)
	if err != nil {
		log.Printf("Failed to find or create user: %v", err)
		http.Error(w, "Failed to process user", http.StatusInternalServerError)
		return
	}

	// Generate JWT token with only user ID
	jwtToken, err := GenerateFamilyCalendarJWT(user.ID)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Render HTML page with token
	renderTokenPage(w, jwtToken, *userInfo)
}

func renderTokenPage(w http.ResponseWriter, token string, userInfo GoogleUserInfo) {
	t, err := template.ParseFiles("auth/templates/auth_success.html")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Token      string
		GivenName  string
		FamilyName string
		Email      string
	}{
		Token:      token,
		GivenName:  userInfo.GivenName,
		FamilyName: userInfo.FamilyName,
		Email:      userInfo.Email,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func generateStateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
