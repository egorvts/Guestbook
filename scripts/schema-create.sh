#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

if [[ -f scripts/.env ]]; then
  source scripts/.env
fi

: "${YDB_CONNECTION_STRING:?Set YDB_CONNECTION_STRING}"

YDB_TOKEN="$(yc iam create-token)"
export YDB_CONNECTION_STRING YDB_TOKEN

echo "Creating/updating the YDB schema..."
cd backend
go run ./cmd --migrate
