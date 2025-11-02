//go:build integration || !ci

package infrastructure_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/domain/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func TestModule(t *testing.T) {
	t.Run("依存関係の注入と初期化", func(t *testing.T) {
		var (
			db           *gorm.DB
			logger       *slog.Logger
			categoryRepo repository.CategoryRepository
			productRepo  repository.ProductRepository
		)

		configOption := fx.Supply(
			fx.Annotate("../../", fx.ResultTags(`name:"configPath"`)),
			fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
		)
		app := fx.New(
			configOption,
			infrastructure.Module,
			fx.Populate(&db, &logger, &categoryRepo, &productRepo),
			fx.NopLogger, // テスト時はログを抑制
		)
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			assert.NoError(t, app.Stop(ctx), "fx app should stop without errors")
		}()

		require.NoError(t, app.Err(), "fx app should initialize without errors")

		// 各依存関係がnilでないことを確認
		assert.NotNil(t, db, "database connection should not be nil")
		assert.NotNil(t, logger, "logger should not be nil")
		assert.NotNil(t, categoryRepo, "category repository should not be nil")
		assert.NotNil(t, productRepo, "product repository should not be nil")

		// データベース接続が有効であることを確認
		conn, err := db.DB()
		require.NoError(t, err, "should get sql.DB from gorm.DB without error")
		assert.NoError(t, conn.Ping(), "should ping database without error")
	})

	t.Run("ライフサイクルフック - アプリケーションの正常起動と停止", func(t *testing.T) {
		var db *gorm.DB

		configOption := fx.Supply(
			fx.Annotate("../../", fx.ResultTags(`name:"configPath"`)),
			fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
		)
		app := fx.New(
			configOption,
			infrastructure.Module,
			fx.Populate(&db),
			fx.NopLogger,
		)
		require.NoError(t, app.Err(), "fx app should initialize without errors")

		// 起動時に接続が有効であることを確認
		conn, err := db.DB()
		require.NoError(t, err, "should get sql.DB from gorm.DB")
		require.NoError(t, conn.Ping(), "ping should succeed after initialization")

		// アプリケーション停止がエラーなく完了することを確認
		// （OnStopフックでDB接続のクローズが正常に実行される）
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = app.Stop(ctx)
		assert.NoError(t, err, "app should stop without errors, including database connection close")
	})

	t.Run("リポジトリの基本動作確認", func(t *testing.T) {
		var (
			categoryRepo repository.CategoryRepository
			productRepo  repository.ProductRepository
		)

		configOption := fx.Supply(
			fx.Annotate("../../", fx.ResultTags(`name:"configPath"`)),
			fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
		)
		app := fx.New(
			configOption,
			infrastructure.Module,
			fx.Populate(&categoryRepo, &productRepo),
			fx.NopLogger,
		)
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			assert.NoError(t, app.Stop(ctx), "fx app should stop without errors")
		}()

		require.NoError(t, app.Err(), "fx app should initialize without errors")

		ctx := context.Background()

		// CategoryRepositoryの動作確認
		categories, err := categoryRepo.List(ctx)
		assert.NoError(t, err, "category repository List should work")
		assert.NotNil(t, categories, "categories should not be nil")

		// ProductRepositoryの動作確認
		products, err := productRepo.List(ctx)
		assert.NoError(t, err, "product repository List should work")
		assert.NotNil(t, products, "products should not be nil")
	})
}
