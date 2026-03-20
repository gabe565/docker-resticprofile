#!/usr/bin/env bash
set -u

: "${CNPG_SECRET_PATH:=/postgresql-app}"

read -r DB_HOST < "$CNPG_SECRET_PATH/host"
read -r DB_NAME < "$CNPG_SECRET_PATH/dbname"
read -r DB_USER < "$CNPG_SECRET_PATH/username"
read -r PGPASSWORD < "$CNPG_SECRET_PATH/password"
export PGPASSWORD

set -x
exec pg_dump --clean --if-exists --no-owner \
  --host="$DB_HOST" --username="$DB_USER" --dbname="$DB_NAME" \
  "$@"
