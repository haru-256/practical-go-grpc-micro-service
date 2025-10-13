.DEFAULT_GOAL := help

.PHONY: ginkgo
ginkgo: ## Run tests
	ginkgo run ./...

.PHONY: ginkgo-detail
ginkgo-detail: ## Run tests with detailed output
	ginkgo run -v ./...

.PHONY: generate-db-models
generate-db-models: ## Generate DB models
	# TODO: dbが立ち上がっていないとエラーになるので、docker-composeで立ち上げるようにする
	cd internal/infrastructure/sqlboiler && \
		sqlboiler mysql -c config/database.toml -o models --no-tests --wipe

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
