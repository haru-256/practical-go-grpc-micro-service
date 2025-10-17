package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
)

type transactionManagerImpl struct{}

func NewTransactionManager() service.TransactionManager {
	return &transactionManagerImpl{}
}

func (tm *transactionManagerImpl) Begin(ctx context.Context) (*sql.Tx, error) {
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return nil, handler.DBErrHandler(err)
	}
	return tx, nil
}

func (tm *transactionManagerImpl) Complete(tx *sql.Tx, err error) error {
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return handler.DBErrHandler(rbErr)
		} else {
			// TODO: ロギングフレームワークを導入したら置き換える。また、DIでloggerを渡すようにする。
			log.Println("トランザクションをロールバックしました")
		}
	} else {
		if cmErr := tx.Commit(); cmErr != nil {
			return handler.DBErrHandler(cmErr)
		} else {
			log.Println("トランザクションをコミットしました")
		}
	}
	return nil
}
