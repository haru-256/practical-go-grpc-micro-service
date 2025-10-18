package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
)

// TransactionManagerImpl はSQLBoilerを使用したTransactionManagerの実装です。
type TransactionManagerImpl struct {
	logger *slog.Logger
}

// NewTransactionManagerImpl は新しいTransactionManagerImplインスタンスを生成します。
// 具象型を返すことで、呼び出し側が必要に応じてインターフェースとして扱えるようにします。
//
// 使用例:
//
//	tm := repository.NewTransactionManagerImpl()
//	var manager service.TransactionManager = tm  // インターフェースとして使用
func NewTransactionManagerImpl(logger *slog.Logger) *TransactionManagerImpl {
	return &TransactionManagerImpl{logger: logger}
}

// Begin は新しいデータベーストランザクションを開始します。
//
// Parameters:
//   - ctx: コンテキスト
//
// Returns:
//   - *sql.Tx: 作成されたトランザクション
//   - error: トランザクション開始時にエラーが発生した場合
func (tm *TransactionManagerImpl) Begin(ctx context.Context) (*sql.Tx, error) {
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return nil, handler.DBErrHandler(err)
	}
	return tx, nil
}

// Complete はトランザクションを完了します。
// errがnilの場合はコミット、非nilの場合はロールバックを実行します。
// コミット/ロールバックの結果はログに記録されます。
//
// Parameters:
//   - ctx: コンテキスト（ログ出力に使用）
//   - tx: 完了するトランザクション
//   - err: トランザクション中に発生したエラー（nilの場合はコミット、非nilの場合はロールバック）
//
// Returns:
//   - error: コミットまたはロールバック時にエラーが発生した場合
func (tm *TransactionManagerImpl) Complete(ctx context.Context, tx *sql.Tx, err error) error {
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return handler.DBErrHandler(rbErr)
		} else {
			tm.logger.WarnContext(ctx, "トランザクションをロールバックしました")
		}
	} else {
		if cmErr := tx.Commit(); cmErr != nil {
			return handler.DBErrHandler(cmErr)
		} else {
			tm.logger.InfoContext(ctx, "トランザクションをコミットしました")
		}
	}
	return nil
}
