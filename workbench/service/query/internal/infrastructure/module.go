package infrastructure

import (
	"context"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/logger"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/domain/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/config"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/db"
	"go.uber.org/fx"
	"gorm.io/gorm"
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
	"infrastructure",
	fx.Provide(
		fx.Annotate(
			config.NewViper,
			fx.ParamTags(`name:"configPath"`, `name:"configName"`),
		),
		db.NewDBConfig,
		db.NewDatabase,
		logger.NewLogger,
		fx.Annotate(
			db.NewCategoryRepositoryImpl,
			fx.As(new(repository.CategoryRepository)),
		),
		fx.Annotate(
			db.NewProductRepositoryImpl,
			fx.As(new(repository.ProductRepository)),
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
func registerLifecycleHooks(lc fx.Lifecycle, db *gorm.DB, logger *slog.Logger) {
	lc.Append(fx.Hook{
		// OnStartはNewDatabaseで接続確認済みなので不要
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "Closing database connection...")
			conn, err := db.DB()
			if err != nil {
				logger.ErrorContext(ctx, "Failed to get database connection", "error", err)
				return err
			}
			err = conn.Close()
			if err != nil {
				logger.ErrorContext(ctx, "Failed to close database connection", "error", err)
				return err
			}
			return nil
		},
	})
}
