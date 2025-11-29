.DEFAULT_GOAL := help

.PHONY: init
init: ## Initial setup
	go mod tidy
	octocov init

.PHONY: lint
lint:  ## Lint proto and go files
	golangci-lint run --config=./.golangci.yml ./...

.PHONY: fmt
fmt:  ## Format proto and go files
	go fmt ./...

.PHONY: test
test: ## Run tests (skip DB tests by default)
	go test $$(go list ./... | grep -v /gen/)

.PHONY: test-ci
test-ci: ## Run tests for CI (skip DB tests)
	go test -tags=ci $$(go list ./... | grep -v /gen/)

.PHONY: test-all
test-all: ## Run all tests including DB tests
	go test -tags=integration $$(go list ./... | grep -v /gen/)

.PHONY: install-tools
install-tools: ## Install tools
	mise install

.PHONY: build-all
build-all: ## Build all services
	docker compose build

.PHONY: up
up: ## up all
	docker compose up -d

.PHONY: up-with-build
up-with-build: ## up all
	docker compose up -d --build

.PHONY: down
down: ## down all
	docker compose down

.PHONY: act-go
act-go: ## Run golang-ci for CI locally
	act push -W=.github/workflows/go.yml --container-architecture linux/amd64

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
