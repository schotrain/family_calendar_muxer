# Family Calendar Muxer

A family calendar application with Google OAuth authentication.

## Prerequisites

- [Docker](https://www.docker.com/) or [OrbStack](https://orbstack.dev/)
- Google OAuth credentials ([Get them here](https://console.cloud.google.com/apis/credentials))

## Development Setup

### 1. Configure Environment Variables

Update the `.env` file in the project root with your Google OAuth credentials:

```bash
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
JWT_SECRET=your-secret-key-change-this-in-production
```

### 2. Start the Application

```bash
docker compose -f docker-compose.dev.yml up --build
```

This will start:
- **Backend** on http://localhost:8080
- **Frontend** on http://localhost:3000
- **PostgreSQL** on port 5432
- **pgAdmin** on http://localhost:5050

### 3. Access the Application

Open your browser and navigate to http://localhost:3000

### 4. Database Management (pgAdmin)

Access the PostgreSQL web UI at http://localhost:5050:
- **Email**: `admin@admin.com`
- **Password**: `admin`

The database connection "Family Calendar DB" is pre-configured and ready to use. Just click on it in the left sidebar to expand and browse the database.

### Stopping Services

```bash
docker compose -f docker-compose.dev.yml down
```
