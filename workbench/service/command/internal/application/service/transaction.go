package service

import (
	"context"
	"database/sql"
)

// TransactionManager はデータベーストランザクションの管理を担うインターフェースです。
// トランザクションの開始、コミット、ロールバックといった操作を提供し、
// アプリケーションサービス層でのトランザクション制御を可能にします。
type TransactionManager interface {
	// Begin は新しいデータベーストランザクションを開始します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//
	// Returns:
	//   - *sql.Tx: 開始されたトランザクション
	//   - error: トランザクション開始に失敗した場合のエラー
	Begin(ctx context.Context) (*sql.Tx, error)

	// Complete はトランザクションを完了します。
	// errがnilの場合はコミット、nilでない場合はロールバックを実行します。
	// この関数はdeferで呼び出されることを想定しています。
	//
	// Parameters:
	//   - tx: 完了するトランザクション
	//   - err: ビジネスロジック実行結果のエラー（nilの場合はコミット、非nilの場合はロールバック）
	//
	// Returns:
	//   - error: コミットまたはロールバックに失敗した場合のエラー
	Complete(tx *sql.Tx, err error) error
}
