VERSION ?= $(shell git describe --tags 2> /dev/null || echo v0)
FILME_BASE ?= $(pwd)

all: binary

.PHONY: binary
binary: ## build binary for Linux
	./scripts/build/binary.sh

.PHONY: run
run: binary
	./bin $(FILME_ARGS)

.PHONY: lint
lint: ## run all the lint tools
	gometalinter ./...
