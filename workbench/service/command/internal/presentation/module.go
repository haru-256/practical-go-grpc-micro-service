package presentation

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	cmdconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/command/v1/commandv1connect"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/presentation/server"
	"go.uber.org/fx"
)

// Module はプレゼンテーション層のFXモジュールです。
// コマンドサービスのプレゼンテーション層に必要な以下の依存関係を提供します：
//   - CategoryServiceHandlerの実装
//   - CommandServerのセットアップ
//   - サーバーライフサイクル管理（起動とグレースフルシャットダウン）
//
// このモジュールは、アプリケーション起動時にgRPCサーバーを開始し、
// アプリケーション停止時にグレースフルシャットダウンを実行するライフサイクルフックを自動的に登録します。
var Module = fx.Module(
	"presentation",
	application.Module,
	fx.Provide(
		fx.Annotate(
			server.NewCategoryServiceHandlerImpl,
			fx.As(new(cmdconnect.CategoryServiceHandler)),
		),
		fx.Annotate(
			server.NewProductServiceHandlerImpl,
			fx.As(new(cmdconnect.ProductServiceHandler)),
		),
		server.NewCommandServer,
	),
	fx.Invoke(registerLifecycleHooks),
)

// registerLifecycleHooks はサーバーのライフサイクルフックをFXアプリケーションライフサイクルに登録します。
// この関数はアプリケーション初期化時にFXフレームワークによって自動的に呼び出されます。
//
// Parameters:
//   - lc: FXライフサイクルマネージャー
//   - srv: 管理対象のコマンドサーバーインスタンス
//   - logger: サーバーライフサイクルイベントを記録するロガー
//
// この関数は以下の2つのフックを登録します：
//   - OnStart: 別のゴルーチンでgRPCサーバーを起動し、起動イベントをログに記録
//   - OnStop: 30秒のタイムアウトでグレースフルシャットダウンを実行し、停止イベントをログに記録
func registerLifecycleHooks(lc fx.Lifecycle, srv *server.CommandServer, logger *slog.Logger) {
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
