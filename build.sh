#!/usr/bin/env bash

set -e -u -o pipefail

source version.sh

echo "Build docker container..."
docker build \
  --build-arg "VERSION=${VERSION}" \
  -t "${IMAGE_NAME}" .
