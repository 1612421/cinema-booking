#!/bin/bash

set -e

# Build and Start containers
echo "🚀 Starting docker containers..."
docker-compose up --build -d

# Wait for MySQL to be ready
echo "⏳ Waiting for MySQL to be ready..."
until docker exec -i cinema_booking_mysql mysqladmin ping -h"localhost" --silent; do
  sleep 2
done

# Run seed.sql
echo "🌱 Seeding database..."
docker exec -i cinema_booking_mysql mysql -uroot -ppassword cinema_booking < ./scripts/database/seed.sql

echo "✅ Setup complete!"
