.DEFAULT_GOAL := help

.PHONY: up
up:  ## Start the database
	docker compose up -d --build

.PHONY: down
down:  ## Stop the database
	docker compose down

.PHONY: dump
dump:  ## Dump the database
	docker compose exec command_db /etc/scripts/dump.sh

.PHONY: restore
restore:  ## Restore the database
	cp ./command/ddl/master.db ./query/ddl/master.db
	docker compose exec query_db \
		bash -c "mysql -uroot -ppassword </etc/ddl/master.db"

.PHONY: start-replication
start-replication:  ## Start replication
	docker compose exec query_db \
		bash -c "mysql -uroot -ppassword </etc/ddl/replication.sql"

.PHONY: create-data
create-data:  ## Create data
	docker compose exec command_db \
		bash -c "mysql -uroot -ppassword </etc/ddl/create_object.sql"
	docker compose exec command_db \
		bash -c "mysql -uroot -ppassword </etc/ddl/create_record.sql"

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
