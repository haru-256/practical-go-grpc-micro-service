package service

import (
	"context"
	"database/sql"
)

// TransactionManager はデータベーストランザクションの管理を担うインターフェースです。
//
//go:generate go tool mockgen -source=$GOFILE -destination=./mock_transaction.go -package=service
type TransactionManager interface {
	// Begin は新しいデータベーストランザクションを開始します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//
	// Returns:
	//   - *sql.Tx: トランザクション
	//   - error: エラー
	Begin(ctx context.Context) (*sql.Tx, error)

	// Complete はトランザクションを完了します。
	// errがnilの場合はコミット、nilでない場合はロールバックを実行します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - err: エラー
	//
	// Returns:
	//   - error: エラー
	Complete(ctx context.Context, tx *sql.Tx, err error) error
}
