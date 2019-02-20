#!/usr/bin/env bash
set -eu -o pipefail

source ./scripts/build/.variables.sh

export CGO_ENABLED=0

go build -o "${TARGET}" --ldflags "${LDFLAGS}" "${SOURCE}"

ln -sf "${TARGET}" bin
