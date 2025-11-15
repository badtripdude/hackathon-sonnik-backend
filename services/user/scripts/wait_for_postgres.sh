#!/bin/sh

set -e

retry() {
  max_attempts="$1"; shift
  seconds="$1"; shift
  cmd="$@"
  attempt_num=1

  until $cmd; do
    if [ "$attempt_num" -ge "$max_attempts" ]; then
      echo "Attempt $attempt_num failed and there are no more attempts left!" >&2
      return 1
    else
      echo "Attempt $attempt_num failed! Trying again in $seconds seconds..." >&2
      attempt_num=$((attempt_num + 1))
      sleep "$seconds"
    fi
  done
}

export PGPASSWORD="${POSTGRES_PASSWORD:-rootroot}"

retry 10 3 psql --host="$POSTGRES_HOST" --port="$POSTGRES_PORT" --username="$POSTGRES_USER" --dbname="$POSTGRES_DB" -c "\\l" >/dev/null 2>&1

echo "$(date +%Y%m%dT%H%M%S) Postgres is up - executing command" >&2

exec "$@"