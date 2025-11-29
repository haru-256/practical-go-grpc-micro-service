.DEFAULT_GOAL := help

.PHONY: init
init: ## Initial setup
	mise install
	buf config init

.PHONY: generate
generate:  ## Generate gRPC code
	buf dep update
	buf generate

.PHONY: lint
lint:  ## Lint proto and go files
	buf lint

.PHONY: fmt
fmt:  ## Format proto and go files
	buf format -w .

.PHONY: curl-grpc
curl-grpc: ## Test with curl (gRPC)
	buf curl --protocol grpc --http2-prior-knowledge \
		--schema . \
		--data '{"sentence": "I am haru256."}' \
		http://localhost:8080/chat.v1.ChatService/Say

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
