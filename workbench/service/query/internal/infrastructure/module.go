package infrastructure

import (
	"context"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/log"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/domain/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/config"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/db"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module はインフラストラクチャ層のFxモジュールです。
// 設定読み込み、データベース接続、リポジトリ実装、ロガーを提供します。
var Module = fx.Module(
	"infrastructure",
	fx.Provide(
		fx.Annotate(
			config.NewViper,
			fx.ParamTags(`name:"configPath"`, `name:"configName"`),
		),
		db.NewDBConfig,
		db.NewDatabase,
		log.NewLogger,
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
//
// Parameters:
//   - lc: Fxライフサイクル
//   - db: データベース接続
//   - logger: ロガー
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
