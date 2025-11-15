package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	cmdconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/command/v1/commandv1connect"
	interceptor "github.com/haru-256/practical-go-grpc-micro-service/pkg/interceptor/connect"
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
func NewCommandServer(viper *viper.Viper, logger *slog.Logger, csh cmdconnect.CategoryServiceHandler, psh cmdconnect.ProductServiceHandler) (*CommandServer, error) {
	reqRespLogger := interceptor.NewReqRespLogger(logger)
	validator, err := interceptor.NewValidator(logger)
	if err != nil {
		return nil, err
	}
	interceptors := connect.WithInterceptors(
		reqRespLogger.NewUnaryInterceptor(), validator.NewUnaryInterceptor(),
	)

	mux := http.NewServeMux()

	// パスとハンドラを登録
	path, handler := cmdconnect.NewCategoryServiceHandler(csh, interceptors)
	mux.Handle(path, handler)
	path, handler = cmdconnect.NewProductServiceHandler(psh, interceptors)
	mux.Handle(path, handler)

	// ヘルスチェックハンドラの登録
	// 標準的なパス: /grpc.health.v1.Health/Check が自動で作られる
	checker := grpchealth.NewStaticChecker(
		"", // 空文字を登録するとサーバー全体のヘルスチェックになる
	)
	mux.Handle(grpchealth.NewHandler(checker))

	// reflection
	reflector := grpcreflect.NewStaticReflector(
		cmdconnect.CategoryServiceName,
		cmdconnect.ProductServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port")),
		Handler:      h2c.NewHandler(mux, &http2.Server{}),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		// HTTP Keep-Aliveのタイムアウト設定
		IdleTimeout: 120 * time.Second,
	}
	return &CommandServer{Server: server}, nil
}
