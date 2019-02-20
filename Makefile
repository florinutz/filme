VERSION ?= $(shell git describe --tags 2> /dev/null || echo v0)
FILME_BASE ?= $(pwd)

.PHONY: binary run lint test help fmt

.DEFAULT_GOAL := help

all: test binary

binary: ## build binary for Linux
	./scripts/build/binary.sh

run: binary
	./bin $(FILME_ARGS)

lint: ## run all the lint tools
	gometalinter ./...

test:
	go test -v -race ./...

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'