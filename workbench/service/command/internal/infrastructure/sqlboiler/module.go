package sqlboiler

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/logger"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/config"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/repository"
	"go.uber.org/fx"
)

// Module はSQLBoilerを使用したインフラストラクチャ層のFxモジュールです。
// このモジュールは以下を提供します:
//   - データベース設定の読み込み（NewDBConfig）
//   - データベース接続の確立（NewDatabase）
//   - カテゴリリポジトリの実装（NewCategoryRepositoryImpl → categories.CategoryRepository）
//   - 商品リポジトリの実装（NewProductRepositoryImpl → products.ProductRepository）
//   - トランザクションマネージャーの実装（NewTransactionManagerImpl → service.TransactionManager）
//   - アプリケーション停止時のDB接続クローズ処理
var Module = fx.Module(
	"sqlboiler",
	config.Module,
	fx.Provide(
		handler.NewDBConfig,
		handler.NewDatabase,
		logger.NewLogger,
		fx.Annotate(
			repository.NewCategoryRepositoryImpl,
			fx.As(new(categories.CategoryRepository)),
		),
		fx.Annotate(
			repository.NewProductRepositoryImpl,
			fx.As(new(products.ProductRepository)),
		),
		fx.Annotate(
			repository.NewTransactionManagerImpl,
			fx.As(new(service.TransactionManager)),
		),
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
func registerLifecycleHooks(lc fx.Lifecycle, db *sql.DB, logger *slog.Logger) {
	lc.Append(fx.Hook{
		// OnStartはNewDatabaseで接続確認済みなので不要
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "Closing database connection...")
			return db.Close()
		},
	})
}
