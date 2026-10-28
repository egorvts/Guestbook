#!/bin/bash
set -e

cd "$(dirname "${BASH_SOURCE[0]}")/../frontend"
BUCKET="${BUCKET:-guestbook.vetoshkin.com}"

echo "Installing dependencies..."
npm ci

echo "Building..."
npm run build

TOKEN="$(yc iam create-token)"
echo "Uploading to YC bucket..."
# Upload assets before the index that references them.
YC_IAM_TOKEN="$TOKEN" yc storage s3 cp dist/ "s3://$BUCKET/" \
  --recursive --exclude 'index.html' --cache-control "no-cache"
YC_IAM_TOKEN="$TOKEN" yc storage s3 cp dist/index.html "s3://$BUCKET/index.html" \
  --cache-control "no-cache"
