#!/bin/bash -e
source version.sh

echo Image: ${IMAGE_NAME}

echo Pushing image...
docker push ${IMAGE_NAME}
