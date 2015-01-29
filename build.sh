#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source version.sh

echo "Build static executable..."
readonly LD_FLAGS="-w -X main.Version=${VERSION}"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags "${LD_FLAGS}" -o "goecho" "github.com/xperimental/goecho"

echo "Build docker container..."
docker build -t "${IMAGE_NAME}" .
