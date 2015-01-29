#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

readonly VERSION=$(git describe --tags --dirty | grep -v dirty || echo "latest")
readonly IMAGE_NAME="xperimental/goecho:${VERSION}"
