package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	queryconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1/queryv1connect"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// QueryServer はgRPCクエリサービスのHTTPサーバーをラップする構造体です。
type QueryServer struct {
	Server *http.Server
}

// NewQueryServer はQueryServerの新しいインスタンスを作成します。
//
// Parameters:
//   - addr: サーバーがリッスンするアドレス (例: ":8080", "localhost:8080")
//   - logger: 構造化ロギング用のslogロガー
//   - csh: カテゴリサービスのgRPCハンドラ実装
//
// Returns:
//   - *QueryServer: 初期化されたQueryServerインスタンス
//   - error: 初期化中にエラーが発生した場合のエラー (現在は常にnil)
func NewQueryServer(viper *viper.Viper, logger *slog.Logger, csh queryconnect.CategoryServiceHandler, psh queryconnect.ProductServiceHandler) (*QueryServer, error) {
	reqRespLogger := NewReqRespLogger(logger)
	validator, err := NewValidator(logger)
	if err != nil {
		return nil, err
	}
	interceptors := connect.WithInterceptors(
		reqRespLogger.NewUnaryInterceptor(),
		validator.NewUnaryInterceptor(),
	)

	mux := http.NewServeMux()

	// パスとハンドラを登録
	path, handler := queryconnect.NewCategoryServiceHandler(csh, interceptors)
	mux.Handle(path, handler)
	path, handler = queryconnect.NewProductServiceHandler(psh, interceptors)
	mux.Handle(path, handler)

	// ヘルスチェックハンドラの登録
	// 標準的なパス: /grpc.health.v1.Health/Check が自動で作られる
	checker := grpchealth.NewStaticChecker(
		queryconnect.CategoryServiceName, // カテゴリサービスを登録
		queryconnect.ProductServiceName,  // プロダクトサービスを登録
	)
	mux.Handle(grpchealth.NewHandler(checker))

	// reflection
	reflector := grpcreflect.NewStaticReflector(
		queryconnect.CategoryServiceName,
		queryconnect.ProductServiceName,
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
	return &QueryServer{Server: server}, nil
}
