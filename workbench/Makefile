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
test: ## Run tests
	go test $$(go list ./... | grep -v /gen/)

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


