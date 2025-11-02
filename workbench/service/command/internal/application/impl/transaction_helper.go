package impl

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
)

// handleTransactionComplete はトランザクションの完了処理を行い、適切なエラーログを出力します。
// この関数はdeferブロック内で使用されることを想定しています。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - tm: トランザクションマネージャー
//   - tx: トランザクション
//   - err: ビジネスロジック実行時のエラー（ポインタで渡すことでerrorを更新可能にする）
//   - result: 処理結果のポインタ（エラー時にゼロ値に設定される）。nilを渡すことも可能（その場合は何もしない）
//   - logger: ログ出力用のロガー
//
// この関数は以下の処理を行います:
//   - err == nil の場合: コミットを試行し、失敗時はエラーログを出力してerrを更新し、resultをゼロ値に設定
//   - err != nil の場合: ロールバックを試行し、失敗時は元のエラーとロールバックエラーの両方をログ出力し、resultをゼロ値に設定
func handleTransactionComplete[T any](ctx context.Context, tm service.TransactionManager, tx *sql.Tx, err *error, result *T, logger *slog.Logger) {
	completeErr := tm.Complete(ctx, tx, *err)
	if *err == nil {
		if completeErr != nil {
			*err = completeErr
			if result != nil {
				var zero T
				*result = zero
			}
			logger.ErrorContext(ctx, "トランザクションのコミットに失敗しました", slog.Any("error", *err))
		}
	} else {
		if result != nil {
			var zero T
			*result = zero
		}
		if completeErr != nil {
			logger.ErrorContext(ctx, "トランザクションのロールバックに失敗しました", slog.Any("original_error", *err), slog.Any("rollback_error", completeErr))
		}
	}
}
