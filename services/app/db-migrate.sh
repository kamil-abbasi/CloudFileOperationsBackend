#!/bin/sh
set -eu

MIGRATE_DB_URL="postgresql://${PG_USER:-admin}:${PG_PASSWORD:-admin}@${PG_HOST:-database}/${PG_DB:-app}?sslmode=disable"

migrate -database "${MIGRATE_DB_URL}" -path ./db/migrations up