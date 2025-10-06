# Family Calendar Muxer

A backend service for managing family calendars with Google OAuth authentication and REST API endpoints.

## Setup

### Prerequisites

- Go 1.25.1 or later
- Google OAuth credentials (Client ID and Secret)

### Environment Configuration

Copy the example environment file:

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

### Install Dependencies

```bash
cd backend
go mod download
```

## Running the Server

### Development with Docker Compose

Run the full stack (backend, frontend, PostgreSQL, and pgAdmin):

```bash
docker-compose -f docker-compose.dev.yml up
```

Services:
- **Backend**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **PostgreSQL**: localhost:5432
- **pgAdmin**: http://localhost:5050

#### pgAdmin Access

Access the PostgreSQL web UI at http://localhost:5050:
- **Email**: `admin@admin.com`
- **Password**: `admin`

The database connection "Family Calendar DB" is pre-configured and ready to use. Just click on it in the left sidebar to expand and browse the database.

### Local Development

```bash
cd backend
go run main.go
```

The server will start on `http://localhost:8080`

## Usage

### 1. Log in with Google OAuth

Open your browser and navigate to:

```
http://localhost:8080/auth/google
```

This will:
1. Redirect you to Google's login page
2. After authentication, redirect back to the callback URL
3. Display your JWT token on a success page

Copy the JWT token from the page.

### 2. Make API calls with the token

Use the token in the `Authorization` header for authenticated API calls:

```bash
# Get user info
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/userinfo
```

Example response:
```json
{
  "id": 1,
  "given_name": "John",
  "family_name": "Doe",
  "email": "john@example.com"
}
```

### Public endpoints

Health check (no authentication required):
```bash
curl http://localhost:8080/health
```

## Testing

Run all tests:
```bash
cd backend
go test ./...
```

Run tests with coverage:
```bash
go test ./... -cover
```

Generate HTML coverage report:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```
