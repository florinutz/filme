#!/usr/bin/env bash
set -eu -o pipefail

source ./scripts/build/.variables.sh

echo "Building statically linked $TARGET"
export CGO_ENABLED=0

go build -o "${TARGET}" --ldflags "${LDFLAGS}" "${SOURCE}"

ln -sf "$(basename "${TARGET}")" bin
