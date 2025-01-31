#!/bin/bash
set -e

echo "🚀 Creating additional databases..."

# List of databases to create
DBS=("api_db" "test_db")

for DB in "${DBS[@]}"; do
  echo "📌 Checking if database '$DB' exists..."
  psql -d postgres -U "$POSTGRES_USER" -tc "SELECT 1 FROM pg_database WHERE datname = '$DB'" | grep -q 1 || psql -d postgres -U "$POSTGRES_USER" -c "CREATE DATABASE $DB;"
  echo "✅ Database '$DB' is ready!"
done

echo "✅ All databases created successfully!"
