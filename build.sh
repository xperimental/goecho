#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source version.sh

echo "Build docker container..."
docker build \
  --build-arg "VERSION=${VERSION}" \
  --build-arg "PACKAGE=github.com/xperimental/goecho" \
  -t "${IMAGE_NAME}" .
