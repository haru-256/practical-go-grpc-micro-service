.DEFAULT_GOAL := help

.PHONY: act-go
act-go: ## Run golang-ci for CI locally
	act push -W=.github/workflows/go.yml

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
