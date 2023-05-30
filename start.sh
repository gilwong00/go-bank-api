#!/bin/sh

#  script exists immediately if returns a non 0 status
set -e

echo "run db migrations"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start app"
# takes all parameters passed to the script and run it
exec "$@"