#!/bin/bash

set -euo pipefail

IMAGE_NAME="karlos98/gus-client"
TAG="latest"
PLATFORMS="linux/amd64,linux/arm64"

echo "✅ Building multi-arch Docker image: $IMAGE_NAME:$TAG"
docker buildx create --use --name gus-builder || docker buildx use gus-builder

docker buildx build \
  --platform "${PLATFORMS}" \
  -t "${IMAGE_NAME}:${TAG}" \
  --push \
  .

echo "🎉 Done! Image pushed to Docker Hub: ${IMAGE_NAME}:${TAG}"

echo "🔍 Available platforms:"
docker buildx imagetools inspect "${IMAGE_NAME}:${TAG}" | grep Platform
