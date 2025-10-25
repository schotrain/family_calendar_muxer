# Family Calendar Muxer - Backend

A backend service for managing family calendars with Google OAuth authentication and REST API endpoints.

## Prerequisites

- Go 1.25.1 or later
- PostgreSQL or SQLite (for database)
- Google OAuth credentials (Client ID and Secret)

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Configure Environment

Copy the example environment file and configure appropriately:

```bash
cp .env.example .env

# Edit .env with your configuration
# See .env.example for required variables
```

## Running the Server

Start the backend server:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `GET /auth/google` - Initiate Google OAuth flow
- `GET /auth/google/callback` - OAuth callback handler

### Public Endpoints
- `GET /health` - Health check (no authentication required)

### Protected Endpoints
Require `Authorization: Bearer <token>` header:
- `GET /api/userinfo` - Get current user information
- `GET /api/calendar-mux` - List user's calendar muxes
- `POST /api/calendar-mux` - Create a new calendar mux
- `DELETE /api/calendar-mux/:id` - Delete a calendar mux

## Building

### Local Build

Build the application:

```bash
go build -o server .
```

Run the built binary:

```bash
./server
```

### Docker Production Build

Build the production Docker image:

```bash
docker build -f Dockerfile.prod -t family-calendar-muxer-backend:prod .
```

Run the production container:

```bash
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://user:password@host:5432/dbname" \
  -e GOOGLE_CLIENT_ID="your-google-client-id" \
  -e GOOGLE_CLIENT_SECRET="your-google-client-secret" \
  -e GOOGLE_REDIRECT_URL="http://localhost:8080/auth/google/callback" \
  -e JWT_SECRET="your-jwt-secret" \
  -e CORS_ALLOWED_ORIGIN="http://localhost:3000" \
  family-calendar-backend:prod
```

**Note:** The production image uses PostgreSQL only (SQLite is not included). The binary is statically compiled for optimal performance and security.

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test ./... -cover
```

Run tests for a specific package:
```bash
go test ./db/services -v
go test ./rest_api_handlers -v
```

Generate HTML coverage report:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```
