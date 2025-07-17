#!/bin/sh
echo "⏳ Waiting for MySQL to be ready..."

until nc -z -v -w30 db 3306
do
  echo "Waiting for database connection..."
  sleep 2
done

echo "✅ MySQL is up. Starting application..."
exec "$@"
