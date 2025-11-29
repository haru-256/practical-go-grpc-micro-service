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

.PHONY: list-methods
list-methods: ## Run buf curl commands
	buf curl http://localhost:8085 --list-methods  --http2-prior-knowledge

.PHONY: curl-list-valid
curl-list-valid: ## Run buf curl commands
	buf curl --protocol grpc --http2-prior-knowledge \
		--schema ../../api \
		--data '{}' \
		http://localhost:8085/query.v1.CategoryService/ListCategories

.PHONY: curl-get-invalid
curl-get-invalid: ## Run buf curl commands
	buf curl --protocol grpc --http2-prior-knowledge \
		--schema ../../api \
		--data '{"id": "hoge"}' \
		http://localhost:8085/query.v1.CategoryService/GetCategoryById

.PHONY: curl-search-valid
curl-search-valid: ## Run buf curl commands
	buf curl --protocol grpc --http2-prior-knowledge \
		--schema ../../api \
		--data '{"keyword": "ペン"}' \
		http://localhost:8085/query.v1.ProductService/SearchProductsByKeyword

.PHONY: run-server
run-server: ## Run query service server
	go run ./cmd/server/main.go

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
