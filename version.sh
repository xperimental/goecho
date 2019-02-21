#!/usr/bin/env bash

set -e -u -o pipefail

readonly VERSION=$(git describe --tags --dirty | grep -v dirty || echo "latest")
readonly IMAGE_NAME="xperimental/goecho:${VERSION}"
