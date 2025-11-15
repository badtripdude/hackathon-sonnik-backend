#!/bin/sh
set -e

export PGPASSWORD="${POSTGRES_PASSWORD:-rootroot}"

echo "$(date +%Y%m%dT%H%M%S) Waiting for users table in Postgres..." >&2

until psql --host="$POSTGRES_HOST" --port="$POSTGRES_PORT" --username="$POSTGRES_USER" --dbname="$POSTGRES_DB" -c "SELECT 1 FROM users LIMIT 1;" >/dev/null 2>&1; do
  echo "Users table not ready, waiting 2 seconds..." >&2
  sleep 2
done

echo "Users table exists, ready to start app" >&2

exec ./mock
