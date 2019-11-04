VERSION ?= $(shell git describe --tags 2> /dev/null || echo v0)
FILME_BASE ?= $(pwd)

.PHONY: binary run lint test help fmt update_test_data

.DEFAULT_GOAL := help

all: test binary ## tests and builds the binary. can be used in ci

binary: ## build binary for Linux
	./scripts/build/binary.sh

run: binary ## builds and runs
	./bin $(ARGS)

lint: ## run all the lint tools
	golint ./...

test: ## run all tests
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

update_test_data: update_test_data_coll33tx update_test_data_google update_test_data_imdb ## does what it says

update_test_data_imdb: ## updates imdb test data
	@go run pkg/collector/imdb/html/update.go pkg/collector/imdb/test-data

update_test_data_coll33tx: ## updates 1337x test data
	@go run pkg/collector/coll33tx/html/update.go pkg/collector/coll33tx/test-data

update_test_data_google: ## updates google test data
	@go run pkg/collector/google/html/update.go pkg/collector/google/test-data

TF_VAR_bucket=lambda-example-flo
TF_VAR_file=build.zip
TF_VAR_ver=v1.0.0

lambda_up: ## creates lambda function
	zip build.zip lambda/main.js
	aws s3 cp build.zip s3://${TF_VAR_bucket}/${TF_VAR_ver}/${TF_VAR_file}
	cd lambda; terraform apply -auto-approve

lambda_down: ## runs terraform destroy for the lambda func
	cd lambda; terraform destroy -auto-approve

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
