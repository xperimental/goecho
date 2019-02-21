#!/usr/bin/env bash

set -e -u -o pipefail

source version.sh

echo Image: ${IMAGE_NAME}

echo Pushing image...
docker push ${IMAGE_NAME}
