#!/bin/sh
# Script to inject runtime environment variables into the built React app

set -e

echo "=== Starting environment variable injection ==="
echo "PUBLIC_AUTH_LOGIN_URL: ${PUBLIC_AUTH_LOGIN_URL:-NOT SET}"
echo "PUBLIC_API_BASE_URL: ${PUBLIC_API_BASE_URL:-NOT SET}"

# Create a JavaScript file that will be loaded by the app
cat <<EOF > /usr/share/nginx/html/env-config.js
window.ENV = {
  AUTH_LOGIN_URL: '${PUBLIC_AUTH_LOGIN_URL:-http://localhost:8080/auth/google}',
  API_BASE_URL: '${PUBLIC_API_BASE_URL:-http://localhost:8080}'
};
EOF

echo "=== Environment configuration injected ==="
cat /usr/share/nginx/html/env-config.js
echo "=== Injection complete ==="

exit 0
