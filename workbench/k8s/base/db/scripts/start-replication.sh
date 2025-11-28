#!/usr/bin/env bash

set -e
set -o pipefail

# 環境変数チェック
REQUIRED_VARS=("COMMAND_DB_HOST" "QUERY_DB_HOST" "MYSQL_ROOT_PASSWORD" "COMMAND_DB_PORT")
for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Error: Environment variable '$var' is required but not set."
        exit 1
    fi
done

# NOTE: psコマンド等でパスワードが見えないようにするため、MYSQL_PWD環境変数を使用
# mysql コマンドラインクライアントは、この環境変数が設定されている場合、自動的にパスワードとして使用する
export MYSQL_PWD="$MYSQL_ROOT_PASSWORD"

echo "Start replication setup..."

# Query DBをリセット (既存データがある場合に備えて)
echo "Resetting Query DB..."
mysql -h "$QUERY_DB_HOST" -uroot -e "STOP REPLICA; RESET REPLICA ALL; RESET MASTER;"

# Dump all databases from Command DB and restore to Query DB
# --set-gtid-purged=ON により、ダンプ内にGTID設定SQLが含まれるため、リストアと同時にQuery DBのGTID_PURGEDがセットされる
echo "Dumping all databases from Command DB and restoring to Query DB..."
mysqldump -h "$COMMAND_DB_HOST" -uroot \
  --single-transaction \
  --set-gtid-purged=ON --routines --triggers \
  --events --hex-blob --all-databases \
  | mysql -h "$QUERY_DB_HOST" -uroot
echo "Dump and Restore completed via pipe."

# Start replication on Query DB
# 既にGTID_PURGEDはセットされているため、MASTER_AUTO_POSITION=1 だけでOK
echo "Starting replication on Query DB..."
# FIXME: rootユーザーは特権を持っており、replication用ユーザーとして適切ではないため、専用ユーザーを作成して使用すること
mysql -h "$QUERY_DB_HOST" -uroot <<EOF
CHANGE REPLICATION SOURCE TO
  SOURCE_HOST='$COMMAND_DB_HOST',
  SOURCE_PORT=$COMMAND_DB_PORT,
  SOURCE_USER='root',
  SOURCE_PASSWORD='$MYSQL_ROOT_PASSWORD',
  SOURCE_AUTO_POSITION=1;
START REPLICA;
EOF

echo "Replication started successfully."
