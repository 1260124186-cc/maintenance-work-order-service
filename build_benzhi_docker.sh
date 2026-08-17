#!/usr/bin/env sh
set -eu

platform="${1:-linux/arm64}"
image="maintenance-work-order-service:${platform#linux/}"
container_name="maintenance-work-order-service-check"

docker buildx build --platform "$platform" --load -t "$image" -f benzhi.Dockerfile .
container_id="$(docker run -d --rm --name "$container_name" -p 8080:8080 "$image")"
cleanup() {
  docker stop "$container_id" >/dev/null 2>&1 || true
}
trap cleanup EXIT INT TERM

docker exec "$container_id" go build ./...
until curl --fail --silent http://127.0.0.1:8080/health >/dev/null; do
  sleep 1
done
curl --fail --silent http://127.0.0.1:8080/assets >/dev/null
