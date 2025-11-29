#!/usr/bin/env bash

mysqldump -uroot -ppassword --single-transaction \
  --set-gtid-purged=COMMENTED --routines --triggers \
  --events --hex-blob --all-databases >/etc/dump/master.db
