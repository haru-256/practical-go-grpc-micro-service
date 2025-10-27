package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	cmdconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/command/v1/commandv1connect"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// CommandServer はgRPCコマンドサービスのHTTPサーバーをラップする構造体です。
type CommandServer struct {
	Server *http.Server
}

// NewCommandServer はCommandServerの新しいインスタンスを作成します。
//
// Parameters:
//   - addr: サーバーがリッスンするアドレス (例: ":8080", "localhost:8080")
//   - logger: 構造化ロギング用のslogロガー
//   - csh: カテゴリサービスのgRPCハンドラ実装
//
// Returns:
//   - *CommandServer: 初期化されたCommandServerインスタンス
//   - error: 初期化中にエラーが発生した場合のエラー (現在は常にnil)
func NewCommandServer(viper *viper.Viper, logger *slog.Logger, csh cmdconnect.CategoryServiceHandler) (*CommandServer, error) {
	serverLogger := NewServerLogger(logger)
	interceptors := connect.WithInterceptors(serverLogger.NewUnaryInterceptor())

	mux := http.NewServeMux()

	// パスとハンドラを登録
	path, handler := cmdconnect.NewCategoryServiceHandler(csh, interceptors)
	mux.Handle(path, handler)

	// reflection
	reflector := grpcreflect.NewStaticReflector(
		cmdconnect.CategoryServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port")),
		Handler:      h2c.NewHandler(mux, &http2.Server{}),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return &CommandServer{Server: server}, nil
}
