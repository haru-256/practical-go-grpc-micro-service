//go:build integration || !ci

package presentation_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"buf.build/gen/go/grpc/grpc/connectrpc/go/grpc/health/v1/healthv1connect"
	healthv1 "buf.build/gen/go/grpc/grpc/protocolbuffers/go/grpc/health/v1"
	"connectrpc.com/connect"
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
		// Arrange
		var (
			categoryServiceHandler queryconnect.CategoryServiceHandler
			productServiceHandler  queryconnect.ProductServiceHandler
			queryServer            *server.QueryServer
		)

		// Act
		app := fx.New(
			configOption,
			presentation.Module,
			fx.Populate(&categoryServiceHandler, &productServiceHandler, &queryServer),
			fx.NopLogger,
		)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			assert.NoError(t, app.Stop(ctx), "fx app should stop without errors")
		})

		// Assert
		require.NoError(t, app.Err(), "fx app should initialize without errors")
		assert.NotNil(t, categoryServiceHandler, "category service handler should not be nil")
		assert.IsType(t, &server.CategoryServiceHandlerImpl{}, categoryServiceHandler, "category service handler should be the correct type")
		assert.NotNil(t, productServiceHandler, "product service handler should not be nil")
		assert.IsType(t, &server.ProductServiceHandlerImpl{}, productServiceHandler, "product service handler should be the correct type")
		assert.NotNil(t, queryServer, "query server should not be nil")
	})

	t.Run("ヘルスチェックエンドポイント", func(t *testing.T) {
		// Arrange
		var queryServer *server.QueryServer

		app := fx.New(
			configOption,
			presentation.Module,
			fx.Populate(&queryServer),
			fx.NopLogger,
		)
		require.NoError(t, app.Err(), "fx app should initialize without errors")
		require.NotNil(t, queryServer, "query server should not be nil")

		// アプリケーションの起動
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		t.Cleanup(cancel)
		require.NoError(t, app.Start(ctx), "fx app should start without errors")
		t.Cleanup(func() {
			stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer stopCancel()
			assert.NoError(t, app.Stop(stopCtx), "fx app should stop without errors")
		})

		// サーバーが初期化されていることを確認
		assert.NotEmpty(t, queryServer.Server.Addr, "server address should be set")

		// Act & Assert: ヘルスチェック
		client := healthv1connect.NewHealthClient(
			http.DefaultClient,
			"http://"+queryServer.Server.Addr,
			connect.WithGRPC(),
		)

		// require.Eventuallyでリトライしながらヘルスチェック
		require.Eventually(t, func() bool {
			clientCtx, clientCancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer clientCancel()

			res, err := client.Check(clientCtx, connect.NewRequest(&healthv1.HealthCheckRequest{}))
			if err != nil {
				return false
			}
			if res == nil || res.Msg == nil {
				return false
			}
			return res.Msg.Status == healthv1.HealthCheckResponse_SERVING
		}, 3*time.Second, 100*time.Millisecond, "service should eventually be in SERVING status")
	})
}
