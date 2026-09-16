#!/usr/bin/env bash
set -euo pipefail

: "${IMAGE_TAG:?IMAGE_TAG is required}"
: "${GHCR_ACTOR:?GHCR_ACTOR is required}"

cd /opt/tiktok-bot

docker login ghcr.io -u "$GHCR_ACTOR" --password-stdin
trap 'docker logout ghcr.io || true' EXIT

export IMAGE_TAG
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d --remove-orphans
printf 'IMAGE_TAG=%s\n' "$IMAGE_TAG" > .deployed-tag
docker image prune -af --filter until=720h
docker compose -f docker-compose.prod.yml ps
