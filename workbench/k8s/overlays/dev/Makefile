.DEFAULT_GOAL := help

.PHONY: db-dry-run
db-dry-run: ## Dry-run DB k8s resources
	kustomize build db/ --load-restrictor LoadRestrictionsNone \
		| kubectl apply -f - --dry-run=server

.PHONY: db-apply
db-apply: ## Apply DB k8s resources
	kustomize build db/ --load-restrictor LoadRestrictionsNone \
		| kubectl apply -f -

.PHONY: services
services: ## Apply services k8s resources
	kustomize build services/ --load-restrictor LoadRestrictionsNone | kubectl apply -f - --dry-run=server

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
