#!/bin/sh
# wait-for-it.sh

set -e

host="$1"
shift
port="$1"
shift
cmd="$@"

until PGPASSWORD=$DB_PASSWORD psql -h "$host" -p "$port" -U "$DB_USER" -d "$DB_NAME" -c '\q' 2>/dev/null; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd