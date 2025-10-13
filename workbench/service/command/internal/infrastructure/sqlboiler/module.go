package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/repository"
	"go.uber.org/fx"
)

// Module はSQLBoilerを使用したインフラストラクチャ層のFxモジュールです。
// このモジュールは以下を提供します:
//   - データベース設定の読み込み（NewDBConfig）
//   - データベース接続の確立（NewDatabase）
//   - カテゴリリポジトリの実装（NewCategoryRepository）
//   - アプリケーション停止時のDB接続クローズ処理
var Module = fx.Module(
	"sqlboiler",
	fx.Provide(
		handler.NewDBConfig,
		handler.NewDatabase,
		repository.NewCategoryRepository,
		repository.NewProductRepository,
	),
	fx.Invoke(registerLifecycleHooks),
)

// registerLifecycleHooks はアプリケーションライフサイクルフックを登録します。
// OnStopフックでデータベース接続のクローズ処理を実行します。
//
// Note: OnStartフックは不要です。NewDatabase関数で既に接続確認が完了しています。
//
// Parameters:
//   - lc: Fxライフサイクル
//   - db: データベース接続
func registerLifecycleHooks(lc fx.Lifecycle, db *sql.DB) {
	lc.Append(fx.Hook{
		// OnStartはNewDatabaseで接続確認済みなので不要
		OnStop: func(ctx context.Context) error {
			// FIXME: 外部からloggerを受け取るようにする(DI)
			log.Println("Closing database connection...")
			return db.Close()
		},
	})
}
