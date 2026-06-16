#!/bin/bash
set -e

IMAGE_NAME="golang-proxy"
CONTAINER_NAME="golang-proxy-app"
DOCKER_USERNAME="${DOCKER_USERNAME:-jefriherditriyanto}"

# Load env vars from .env if present
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

LISTEN_PORT="${LISTEN_PORT:-8080}"
SERVER_PORT="${SERVER_PORT:-8000}"

echo "=== Golang Proxy Docker ==="
echo ""
echo "1) Build & Run locally"
echo "2) Build multi-arch & Push to Docker Hub"
read -p "Choice [1/2]: " choice

if [ "$choice" = "2" ]; then
  DOCKER_HUB_REPO="$DOCKER_USERNAME/$IMAGE_NAME"

  echo "Fetching latest tag from Docker Hub..."
  latest_version=$(curl -s "https://hub.docker.com/v2/repositories/${DOCKER_HUB_REPO}/tags/?page_size=100" | \
    jq -r '.results[].name' 2>/dev/null | \
    grep -E '^[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -n 1)
  latest_version="${latest_version:-0.0.0}"
  echo "Latest: $latest_version"

  while true; do
    read -p "Version tag (x.x.x): " version_tag
    if [[ ! "$version_tag" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
      echo "Invalid format. Use x.x.x"
      continue
    fi
    if [ "$(printf '%s\n' "$latest_version" "$version_tag" | sort -V | tail -n 1)" = "$version_tag" ] && [ "$version_tag" != "$latest_version" ]; then
      break
    fi
    echo "Version must be > $latest_version"
  done

  echo "Building & pushing $DOCKER_HUB_REPO:$version_tag ..."
  docker buildx build --no-cache --platform linux/amd64,linux/arm64 \
    -t "$DOCKER_HUB_REPO:latest" \
    -t "$DOCKER_HUB_REPO:$version_tag" \
    . --push
  echo "Done: $DOCKER_HUB_REPO:latest, $DOCKER_HUB_REPO:$version_tag"

else
  echo "Building Docker image..."
  docker build -t "$IMAGE_NAME" .

  docker stop "$CONTAINER_NAME" 2>/dev/null || true
  docker rm "$CONTAINER_NAME" 2>/dev/null || true

  echo "Starting container..."
  docker run -d \
    -p 80:80 \
    -p 443:443 \
    -p "$LISTEN_PORT:8080" \
    -p "$SERVER_PORT:8000" \
    -e LISTEN_PORT="$LISTEN_PORT" \
    -e SERVER_PORT="$SERVER_PORT" \
    -e SERVER_USERNAME="${SERVER_USERNAME:-admin}" \
    -e SERVER_PASSWORD="${SERVER_PASSWORD:-admin}" \
    --name "$CONTAINER_NAME" \
    "$IMAGE_NAME"

  sleep 2
  docker logs "$CONTAINER_NAME"
  echo ""
  echo "Dashboard: http://localhost:$SERVER_PORT (admin/admin)"
  echo "Proxy    : http://localhost:$LISTEN_PORT"
fi
