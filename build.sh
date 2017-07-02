#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source version.sh

echo "Build docker container..."
docker build -t "${IMAGE_NAME}" .
