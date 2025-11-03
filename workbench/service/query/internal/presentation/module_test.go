//go:build integration || !ci

package presentation_test

import (
	"context"
	"testing"
	"time"

	queryconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1/queryv1connect"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/presentation"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/presentation/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	configOption := fx.Supply(
		fx.Annotate("../../", fx.ResultTags(`name:"configPath"`)),
		fx.Annotate("config", fx.ResultTags(`name:"configName"`)),
	)
	t.Run("依存関係の注入と初期化", func(t *testing.T) {
		var (
			categoryServiceHandler queryconnect.CategoryServiceHandler
			productServiceHandler  queryconnect.ProductServiceHandler
		)

		app := fx.New(
			configOption,
			presentation.Module,
			fx.Populate(&categoryServiceHandler, &productServiceHandler),
			fx.NopLogger, // テスト時はログを抑制
		)
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			assert.NoError(t, app.Stop(ctx), "fx app should stop without errors")
		}()

		require.NoError(t, app.Err(), "fx app should initialize without errors")

		// 各依存関係がnilでないことを確認
		assert.NotNil(t, categoryServiceHandler, "category service handler should not be nil")
		assert.NotNil(t, productServiceHandler, "product service handler should not be nil")
	})

	t.Run("サーバーのライフサイクル", func(t *testing.T) {
		var srv *server.QueryServer

		app := fx.New(
			configOption,
			presentation.Module,
			fx.Populate(&srv),
			fx.NopLogger,
		)

		require.NoError(t, app.Err(), "fx app should initialize without errors")
		assert.NotNil(t, srv, "query server should not be nil")
		assert.NotNil(t, srv.Server, "http server should not be nil")

		// アプリケーションの起動
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		require.NoError(t, app.Start(ctx), "fx app should start without errors")

		// サーバーが初期化されていることを確認
		assert.NotEmpty(t, srv.Server.Addr, "server address should be set")

		// アプリケーションの停止
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()
		assert.NoError(t, app.Stop(stopCtx), "fx app should stop without errors")
	})
}
