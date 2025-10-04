package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// GoogleUserInfo represents the user information from Google
type GoogleUserInfo struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

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
	token, err := GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get user info from Google
	client := GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("Failed to decode user info: %v", err)
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	jwtToken, err := GenerateFamilyCalendarJWT(userInfo.Email, userInfo.GivenName, userInfo.FamilyName)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Render HTML page with token
	renderTokenPage(w, jwtToken, userInfo)
}

func renderTokenPage(w http.ResponseWriter, token string, userInfo GoogleUserInfo) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Authentication Successful</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
        }
        .user-info {
            margin: 20px 0;
            padding: 15px;
            background-color: #f8f9fa;
            border-radius: 4px;
        }
        .token-container {
            margin: 20px 0;
        }
        .token {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 4px;
            word-wrap: break-word;
            font-family: monospace;
            font-size: 12px;
            border: 1px solid #dee2e6;
        }
        button {
            background-color: #007bff;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        button:hover {
            background-color: #0056b3;
        }
        .success {
            color: #28a745;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="success">✓ Authentication Successful</h1>

        <div class="user-info">
            <h3>Welcome, {{.GivenName}} {{.FamilyName}}!</h3>
            <p><strong>Email:</strong> {{.Email}}</p>
        </div>

        <div class="token-container">
            <h3>Your Family Calendar JWT Token:</h3>
            <div class="token" id="token">{{.Token}}</div>
            <button onclick="copyToken()">Copy Token</button>
        </div>

        <p><small>This token is valid for 24 hours. Keep it secure and do not share it.</small></p>
    </div>

    <script>
        function copyToken() {
            const token = document.getElementById('token').textContent;
            navigator.clipboard.writeText(token).then(() => {
                alert('Token copied to clipboard!');
            }).catch(err => {
                console.error('Failed to copy:', err);
            });
        }
    </script>
</body>
</html>
`

	t, err := template.New("token").Parse(tmpl)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
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
