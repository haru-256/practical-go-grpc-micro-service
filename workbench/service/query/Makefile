.DEFAULT_GOAL := help

.PHONY: lint
lint:  ## Lint go files
	$(MAKE) -C ../../ lint

.PHONY: fmt
fmt:  ## Format go files
	$(MAKE) -C ../../ fmt

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test-ci
test-ci: ## Run tests for CI (skip DB tests)
	go test -tags=ci -v ./...

.PHONY: test-all
test-all: ## Run all tests (integration tests included)
	go test -tags=integration -v ./...

.PHONY: generate-mocks
generate-mocks: ## Generate mocks
	go generate ./...

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
