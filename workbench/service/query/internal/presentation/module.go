package presentation

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	queryconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1/queryv1connect"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/presentation/server"
	"go.uber.org/fx"
)

// Module はプレゼンテーション層のFXモジュールです。
var Module = fx.Module(
	"presentation",
	infrastructure.Module,
	fx.Provide(
		server.NewValidator,
		fx.Annotate(
			server.NewCategoryServiceHandlerImpl,
			fx.As(new(queryconnect.CategoryServiceHandler)),
		),
		fx.Annotate(
			server.NewProductServiceHandlerImpl,
			fx.As(new(queryconnect.ProductServiceHandler)),
		),
		server.NewQueryServer,
	),
	fx.Invoke(registerLifecycleHooks),
)

// registerLifecycleHooks はサーバーのライフサイクルフックをFXアプリケーションライフサイクルに登録します。
//
// Parameters:
//   - lc: FXライフサイクルマネージャー
//   - srv: 管理対象のクエリサーバーインスタンス
//   - logger: サーバーライフサイクルイベントを記録するロガー
//
// この関数は以下の2つのフックを登録します：
//   - OnStart: 別のゴルーチンでgRPCサーバーを起動し、起動イベントをログに記録
//   - OnStop: 30秒のタイムアウトでグレースフルシャットダウンを実行し、停止イベントをログに記録
func registerLifecycleHooks(lc fx.Lifecycle, srv *server.QueryServer, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// port競合によるエラーになる問題の対策のため、ListenAndServeではなく、ListenしてServeする
			ln, err := net.Listen("tcp", srv.Server.Addr)
			if err != nil {
				return err
			}
			// サーバーを別のゴルーチンで起動
			go func() {
				logger.Info("Starting Query gRPC server", slog.String("addr", srv.Server.Addr))
				if serveErr := srv.Server.Serve(ln); serveErr != nil && serveErr != http.ErrServerClosed {
					logger.Error("Server failed to start", slog.Any("error", serveErr))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			// サーバーをグレースフルシャットダウン
			if shutdownErr := srv.Server.Shutdown(shutdownCtx); shutdownErr != nil {
				logger.Error("Server shutdown failed", "error", shutdownErr)
				return shutdownErr
			}
			logger.Info("Server stopped gracefully")
			return nil
		},
	})
}
