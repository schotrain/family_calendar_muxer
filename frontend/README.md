# Family Calendar Muxer - Frontend

React TypeScript frontend for the Family Calendar Muxer application built with Rsbuild and Ant Design.

## Setup

1. Install dependencies:

```bash
npm install
```

2. Configure environment variables:

Copy the template file and update with your settings:

```bash
cp .env.template .env
```

Edit `.env` and configure:
- `PUBLIC_AUTH_LOGIN_URL` - The authentication login endpoint (default: `http://localhost:8080/auth/google`)
- `PUBLIC_API_BASE_URL` - The API base URL (default: `http://localhost:8080`)

## Running the App

Start the development server:

```bash
npm run dev
```

The app will open at `http://localhost:3000`

## Building for Production

```bash
npm run build
```

## Preview Production Build

```bash
npm run preview
```
