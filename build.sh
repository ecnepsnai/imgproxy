#!/bin/bash
set -e

VERSION=${1:?Version required}
REVISION=$(git rev-parse HEAD)
DATETIME=$(date --rfc-3339=seconds)
BUILD_DATE=$(date -R)

rm -f docker/imgproxy ./imgproxy
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X 'main.Version=${VERSION}' -X 'main.BuildDate=${BUILD_DATE}'" -trimpath -buildmode=exe -o imgproxy
mv ./imgproxy docker/
cd docker/
podman build \
    --squash \
    --no-cache \
    --format docker \
    --label "org.opencontainers.image.created=${DATETIME}" \
    --label "org.opencontainers.image.version=${VERSION}" \
    --label "org.opencontainers.image.revision=${REVISION}" \
    -t ghcr.io/ecnepsnai/imgproxy:latest \
    -t ghcr.io/ecnepsnai/imgproxy:${VERSION} \
    .
rm -f imgproxy
