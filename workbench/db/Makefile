.DEFAULT_GOAL := help

.PHONY: up
up:  ## Start the database
	docker compose up -d --build

.PHONY: down
down:  ## Stop the database
	docker compose down

.PHONY: dump
dump:  ## command_dbからダンプを取得
	docker compose exec command_db /etc/scripts/dump.sh

.PHONY: restore
restore:  ##  query_dbにdumpしたcommand_dbのデータをリストア
	cp ./command/ddl/master.db ./query/ddl/master.db
	docker compose exec query_db \
		bash -c "mysql -uroot -ppassword </etc/ddl/master.db"

.PHONY: start-replication
start-replication:  ## query_dbでレプリケーションを開始
	docker compose exec query_db \
		bash -c "mysql -uroot -ppassword </etc/ddl/start_replication.sql"

.PHONY: reset-replication
reset-replication:  ## query_dbでレプリケーションをリセット
	docker compose exec query_db \
		bash -c "mysql -uroot -ppassword </etc/ddl/reset_replication.sql"

.PHONY: fetch-gtid-from-dump
fetch-gtid-from-dump:  ## query_dbでdumpしたmaster.dbからGTIDを取得
	docker compose exec query_db /etc/scripts/fetch_gtid.sh

.PHONY: dump-restore
dump-restore:  ## query_dbでダンプを取得してリストア
	$(MAKE) dump
	$(MAKE) restore
	$(MAKE) fetch-gtid-from-dump

.PHONY: check-replication
check-replication:  ## Check replication status
	docker compose exec query_db \
		mysql -uroot -ppassword -e "SHOW REPLICA STATUS\G" \
		| grep -E "(Replica_IO_Running|Replica_SQL_Running|Seconds_Behind_Source)"

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
