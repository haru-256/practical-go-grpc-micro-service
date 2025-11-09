#!/usr/bin/env bash

set -e

# master.dbからGTIDを抽出
GTID_LINE=$(grep "SET @@GLOBAL.GTID_PURGED" /etc/dump/master.db || echo "")

if [ -z "$GTID_LINE" ]; then
    echo "Error: Could not find GTID_PURGED in master.db"
    exit 1
fi

# GTID値を抽出（'...'の中身を取り出す）
GTID=$(echo "$GTID_LINE" | sed -n "s/.*GTID_PURGED='\(.*\)';.*/\1/p")

if [ -z "$GTID" ]; then
    echo "Error: Could not extract GTID value"
    exit 1
fi

echo "Found GTID: $GTID"
