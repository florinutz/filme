#!/usr/bin/env bash
set -eu

VERSION=${VERSION:-"unknown"}
GITCOMMIT=${GITCOMMIT:-$(git rev-parse --short HEAD 2> /dev/null || true)}
BUILDTIME=${BUILDTIME:-$(date --utc --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')}
FILME_BASE=${FILME_BASE:-$(dirname "$(dirname "$(pwd)")")}

export LDFLAGS="\
    -w \
    -X \"github.com/florinutz/filme/pkg.Commit=${GITCOMMIT}\" \
    -X \"github.com/florinutz/filme/pkg.BuildTime=${BUILDTIME}\" \
    -X \"github.com/florinutz/filme/pkg.Version=${VERSION}\" \
    ${LDFLAGS:-} \
"

GOOS="${GOOS:-$(go env GOHOSTOS)}"
GOARCH="${GOARCH:-$(go env GOHOSTARCH)}"

export SOURCE="commands/main/main.go"
export TARGET="build/filme-$GOOS-$GOARCH"
