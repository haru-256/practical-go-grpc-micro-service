package handler

import (
	"errors"
	"log"
	"net"

	"github.com/go-sql-driver/mysql"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

// データベースアクセスエラーのハンドリング
func DBErrHandler(err error) error {
	var opErr *net.OpError
	var driverErr *mysql.MySQLError
	if errors.As(err, &opErr) { // 接続タイムアウトやネットワーク関連の問題で接続が確立できない場合
		// TODO: slogに変更する
		log.Println(err.Error())
		return errs.NewInternalErrorWithCause("DB_CONNECTION_ERROR", opErr.Error(), opErr)
	} else if errors.As(err, &driverErr) { // MySQLドライバエラーの場合
		log.Printf("Code:%d Message:%s", driverErr.Number, driverErr.Message)
		if driverErr.Number == 1062 { // 一意制約違反の場合
			return errs.NewCRUDErrorWithCause("DB_UNIQUE_CONSTRAINT_VIOLATION", "一意制約違反です。", driverErr)
		} else {
			return errs.NewInternalErrorWithCause("DB_DRIVER_ERROR", driverErr.Message, driverErr)
		}
	} else { // その他のエラー
		log.Println(err.Error())
		return errs.NewInternalErrorWithCause("UNKNOWN_ERROR", err.Error(), err)
	}
}
