#!/bin/sh

set -e

echo "Running migrations..."
goose -dir /migrations postgres "$POSTGRES_URL_GOOSE" up

echo "Migrations completed"