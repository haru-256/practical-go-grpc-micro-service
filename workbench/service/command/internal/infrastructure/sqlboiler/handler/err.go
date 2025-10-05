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
	if errors.As(err, &opErr) { // 接続がタイムアウトかネットワーク関連の問題が原因で接続が確立できない?
		// TODO: slogに変更する
		log.Println(err.Error())
		return errs.NewInternalErrorWithCause("DB_CONNECTION_ERROR", opErr.Error(), opErr)
	} else if errors.As(err, &driverErr) { // MySQLドライバエラー?
		log.Printf("Code:%d Message:%s \n", driverErr.Number, driverErr.Message)
		if driverErr.Number == 1062 { // 一意制約違反?
			return errs.NewCRUDErrorWithCause("DB_UNIQUE_CONSTRAINT_VIOLATION", "一意制約違反です。", driverErr)
		} else {
			return errs.NewInternalErrorWithCause("DB_DRIVER_ERROR", driverErr.Message, driverErr)
		}
	} else { // その他のエラー
		log.Println(err.Error())
		return errs.NewInternalErrorWithCause("UNKNOWN_ERROR", err.Error(), err)
	}
}
