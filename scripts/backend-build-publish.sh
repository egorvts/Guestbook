#!/bin/bash
set -eo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

if [[ -f scripts/.env ]]; then
  source scripts/.env
fi

: "${REGISTRY_ID:?Set REGISTRY_ID}"
: "${SERVICE_ACCOUNT_ID:?Set SERVICE_ACCOUNT_ID}"
: "${YDB_CONNECTION_STRING:?Set YDB_CONNECTION_STRING}"
: "${CONTAINER_NAME:?Set CONTAINER_NAME to the existing container used by API Gateway}"

IMAGE_TAG="${IMAGE_TAG:-$(date -u +%Y%m%dT%H%M%SZ)-$RANDOM}"
IMAGE="cr.yandex/$REGISTRY_ID/guestbook:$IMAGE_TAG"

echo "Building..."
docker build --platform linux/amd64 -t "$IMAGE" backend/

echo "Pushing to Container Registry..."
yc container registry configure-docker
docker push "$IMAGE"

echo "Deploying a new revision to $CONTAINER_NAME..."
yc serverless container revision deploy \
  --endpoint api.cloud.yandex.net:443 \
  --container-name "$CONTAINER_NAME" \
  --image "$IMAGE" \
  --service-account-id "$SERVICE_ACCOUNT_ID" \
  --environment "YDB_CONNECTION_STRING=$YDB_CONNECTION_STRING" \
  --cores 1 \
  --memory 256MB \
  --concurrency 10 \
  --execution-timeout 30s

echo "Deployed $IMAGE to $CONTAINER_NAME"
