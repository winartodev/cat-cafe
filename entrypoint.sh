#!/bin/sh
set -e

echo "==================================="
echo "Cat Cafe API - Starting up"
echo "==================================="

# Wait for PostgreSQL
echo "Waiting for PostgreSQL at $DB_HOST:$DB_PORT..."
until nc -z $DB_HOST $DB_PORT; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 2
done
echo "✓ PostgreSQL is up"
echo ""

# Wait for Redis
echo "Waiting for Redis at $REDIS_ADDR..."
REDIS_HOST_ONLY=$(echo $REDIS_ADDR | cut -d':' -f1)
REDIS_PORT_ONLY=$(echo $REDIS_ADDR | cut -d':' -f2)
until nc -z $REDIS_HOST_ONLY ${REDIS_PORT_ONLY:-6379}; do
  echo "Redis is unavailable - sleeping"
  sleep 2
done
echo "✓ Redis is up"

# Run migrations
echo "------------------------------------------"
echo "Running database migrations..."

MIGRATION_PATH="/app/db/migrations"

# Check if migrations folder exists
if [ ! -d "$MIGRATION_PATH" ]; then
    echo "✗ Migration folder $MIGRATION_PATH not found!"
    ls -la /app/
    exit 1
fi

echo ""

# Build DB URL
DB_URL="postgresql://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSL_MODE"

# Run migration with absolute path
migrate -path=$MIGRATION_PATH -database="$DB_URL" -verbose up

if [ $? -eq 0 ]; then
    echo "✓ Migrations completed successfully"
else
    echo "✗ Migration failed"
    exit 1
fi

echo "------------------------------------------"
echo ""

echo "==================================="
echo "Starting application..."
echo "==================================="

exec "$@"
