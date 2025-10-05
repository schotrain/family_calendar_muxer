# Family Calendar Muxer

A backend service for managing family calendars with Google OAuth authentication and REST API endpoints.

## Features

- **Google OAuth/OIDC Authentication** - Secure login flow with JWT token generation
- **REST API** - User management with full CRUD operations
- **SQLite Database** - Persistent data storage with GORM ORM
- **Request/Response Validation** - Schema validation using go-playground/validator
- **Chi Router** - Fast and lightweight HTTP routing

## Prerequisites

- Go 1.25.1 or later
- Google OAuth credentials (Client ID and Secret)

## Getting Started

### 1. Clone and Setup

```bash
git clone <repository-url>
cd family_calendar_muxer
```

### 2. Configure Environment

Copy the example environment file and configure your settings:

```bash
cp backend/.env.example backend/.env
```

Edit `backend/.env` with your Google OAuth credentials:

```env
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
JWT_SECRET=your-secret-key-change-this-in-production
USE_SECURE_CONNECTIONS=false  # Set to true in production with HTTPS
```

### 3. Install Dependencies

```bash
cd backend
go mod download
```

### 4. Run the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Authentication

### Google OAuth Login

Visit `http://localhost:8080/auth/google` to initiate the Google OAuth flow:

1. User is redirected to Google's consent page
2. After approval, Google redirects back to `/auth/google/callback`
3. Backend creates/updates user in database
4. Backend generates a Family Calendar JWT token
5. Token is displayed on a styled HTML page with copy functionality

The JWT token includes:
- **User ID** (local database ID)
- Email address
- Given name and family name
- 24-hour expiration

## API Endpoints

All API endpoints (except `/health`) require authentication via Bearer token.

### Authentication

Include the JWT token in the `Authorization` header:

```bash
Authorization: Bearer <your_jwt_token>
```

### Health Check (Public)

**GET** `/health`

Returns the health status of the server. No authentication required.

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok"
}
```

### Get User Info (Protected)

**GET** `/api/userinfo`

Returns the authenticated user's information based on their JWT token.

**Authentication Required:** Yes

```bash
curl -H "Authorization: Bearer <your_jwt_token>" \
  http://localhost:8080/api/userinfo
```

**Success Response (200):**
```json
{
  "id": 1,
  "given_name": "John",
  "family_name": "Doe",
  "email": "john@example.com"
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "Authorization header required"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "User not found"
}
```

## Project Structure

```
family_calendar_muxer/
├── backend/
│   ├── main.go                          # Server setup and routing
│   ├── auth/
│   │   ├── config.go                    # OAuth and JWT configuration
│   │   ├── jwt.go                       # JWT token generation
│   │   ├── middleware.go                # JWT validation middleware
│   │   ├── handlers.go                  # Auth handlers
│   │   └── templates/
│   │       └── auth_success.html        # Success page template
│   ├── rest_api_handlers/
│   │   ├── health_handler.go            # Health check endpoint
│   │   ├── user_handler.go              # User info endpoint
│   │   ├── user_handler_schema.go       # Response schemas
│   │   └── utils/
│   │       └── response.go              # JSON response helpers
│   ├── db/
│   │   ├── database.go                  # DB connection and migrations
│   │   ├── models/
│   │   │   └── user.go                  # User database model
│   │   └── services/
│   │       └── user_service.go          # User data access
│   ├── go.mod                           # Go dependencies
│   └── .env.example                     # Environment variables template
├── .gitignore
└── README.md
```

## Development

### Database

The application uses SQLite for data persistence. The database file `family_calendar.db` is created automatically on first run in the `backend` directory.

### Environment Variables

- `GOOGLE_CLIENT_ID` - Your Google OAuth Client ID
- `GOOGLE_CLIENT_SECRET` - Your Google OAuth Client Secret
- `GOOGLE_REDIRECT_URL` - OAuth callback URL
- `JWT_SECRET` - Secret key for signing JWT tokens
- `USE_SECURE_CONNECTIONS` - Set to `true` for HTTPS (production), `false` for HTTP (development)

## Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URI: `http://localhost:8080/auth/google/callback`
6. Copy Client ID and Client Secret to `.env` file

## License

[Add your license here]
