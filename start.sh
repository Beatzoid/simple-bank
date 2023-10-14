#!/bin/sh

# Exit immediately if a command exits with a non-zero status
set -e

echo "run db migration"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
