.DEFAULT_GOAL := help

.PHONY: ginkgo
ginkgo: ## Run tests
	go tool ginkgo run ./...

.PHONY: ginkgo-detail
ginkgo-detail: ## Run tests with detailed output
	go tool ginkgo run -v ./...

.PHONY: install-tools
install-tools: ## Install tools
	$(MAKE) -C ../../ install-tools
	# sqlboilerの内部で、われているため、toolディレクティブでは対応できず、PATHに明示的にインストールする必要があるため
	# Error: could not register driver: could not find driver executable: exec: "sqlboiler-mysql": executable file not found in $PATH
	go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-mysql


.PHONY: generate-db-models
generate-db-models: ## Generate DB models
	# TODO: dbが立ち上がっていないとエラーになるので、docker-composeで立ち上げるようにする
	cd internal/infrastructure/sqlboiler && \
		go tool sqlboiler mysql -c config/database.toml -o models --no-tests --wipe

.PHONY: generate-mocks
generate-mocks: ## Generate mocks
	go generate ./...

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
