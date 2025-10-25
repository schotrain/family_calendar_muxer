# Family Calendar Muxer - Frontend

React TypeScript frontend for the Family Calendar Muxer application built with Rsbuild and Ant Design.

## Setup

1. Install dependencies:

```bash
npm install
```

2. Configure environment variables:

Copy the example file and configure appropriately:

```bash
cp .env.example .env

# Edit .env with your configuration
# See .env.example for required variables
```

## Running the App

Start the development server:

```bash
npm run dev
```

The app will open at `http://localhost:3000`

## Building for Production

### Local Build

Build the production bundle:

```bash
npm run build
```

Preview the production build:

```bash
npm run preview
```

### Docker Production Build

Build the production Docker image:

```bash
docker build -f Dockerfile.prod -t family-calendar-muxer-frontend:prod .
```

Run the production container:

```bash
docker run -p 80:80 family-calendar-frontend:prod
```

Or run on a different port:

```bash
docker run -p 3000:80 family-calendar-frontend:prod
```

**Features:**
- Multi-stage build for minimal image size (~50MB or less)
- Nginx server with optimized configuration
- React Router support with client-side routing
- Gzip compression enabled
- Static asset caching (1 year for immutable assets)
- Security headers included
- Health check endpoint configured
