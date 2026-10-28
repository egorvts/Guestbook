#!/bin/bash
set -e

cd "$(dirname "${BASH_SOURCE[0]}")/../frontend"
BUCKET="${BUCKET:-guestbook.vetoshkin.com}"

echo "Installing dependencies..."
npm ci

echo "Building..."
npm run build

TOKEN="$(yc iam create-token)"

upload() {
  local output status=0
  output="$(YC_IAM_TOKEN="$TOKEN" yc storage s3 cp "$@" 2>&1)" || status=$?
  printf '%s\n' "$output"
  # YC can report failed transfers while returning exit code 0.
  if (( status != 0 )) || [[ "$output" == *"upload failed:"* ]]; then
    return 1
  fi
}

echo "Uploading to YC bucket..."
# Upload assets before the index that references them.
upload dist/ "s3://$BUCKET/" \
  --recursive --exclude 'index.html' --cache-control "no-cache"
upload dist/index.html "s3://$BUCKET/index.html" \
  --cache-control "no-cache"
