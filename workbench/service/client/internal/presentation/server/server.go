package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/utils"
	_ "github.com/haru-256/practical-go-grpc-micro-service/service/client/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
)

// CQRSServiceServer はCQRSクライアントサービスのHTTPサーバー
type CQRSServiceServer struct {
	logger *slog.Logger // ロガー
	e      *echo.Echo   // Echoインスタンス
	Addr   string       // サーバーの実際のアドレス（テスト用）
}

// CQRSServiceConfig はCQRSサービスの設定
type CQRSServiceConfig struct {
	Port string // サーバーポート
}

// NewCQRSServiceConfig は設定ファイルからCQRSServiceConfigを生成します。
//
// Parameters:
//   - v: Viperインスタンス
//
// Returns:
//   - *CQRSServiceConfig: 設定のインスタンス
//   - error: 設定の読み込みエラー
func NewCQRSServiceConfig(v *viper.Viper) (*CQRSServiceConfig, error) {
	var configErrors []error
	cfg := &CQRSServiceConfig{
		Port: utils.GetKey[string](v, "server.port", &configErrors),
	}
	if len(configErrors) > 0 {
		return nil, errors.Join(configErrors...)
	}

	return cfg, nil
}

// NewCQRSServiceServer はCQRSServiceServerを生成します。
//
// Parameters:
//   - cfg: サーバー設定
//   - logger: ロガー
//   - handler: HTTPハンドラ
//
// Returns:
//   - *CQRSServiceServer: CQRSServiceServerのインスタンス
func NewCQRSServiceServer(cfg *CQRSServiceConfig, logger *slog.Logger, handler *CQRSServiceHandler) *CQRSServiceServer {
	e := echo.New()
	// Echoのデフォルトロガーを無効化 (二重出力を防ぐため)
	// e.HideBanner = true
	// e.HidePort = true
	// ミドルウェアの設定
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogError:    true,
		HandleError: true, // エラー発生時もログ出力する
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// slogを使ってログ出力
			if v.Error == nil {
				logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("method", v.Method),
				)
			} else {
				logger.LogAttrs(c.Request().Context(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("method", v.Method),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		// ここでslogなどを使ってロギングする
		slog.Info("Response Dump",
			slog.String("method", c.Request().Method),
			slog.String("uri", c.Request().RequestURI),
			slog.Int("status", c.Response().Status),
			slog.String("response_body", string(resBody)), // レスポンスボディを記録
		)
	}))

	// validatorの設定
	e.Validator = NewRequestValidator()

	// ルーティングの設定
	// Swaggerエンドポイントの設定
	// ヘルスチェック用エンドポイント
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// カテゴリ関連のエンドポイント
	e.GET("/categories", handler.CategoryList)
	e.POST("/categories", handler.CreateCategory)
	e.GET("/categories/:id", handler.CategoryById)
	e.PUT("/categories/:id", handler.UpdateCategory)
	e.DELETE("/categories/:id", handler.DeleteCategory)

	// 商品関連のエンドポイント
	e.GET("/products", handler.ProductList) // keywordパラメータがある場合は検索、ない場合は一覧取得
	e.POST("/products", handler.CreateProduct)
	e.GET("/products/:id", handler.ProductById)
	e.PUT("/products/:id", handler.UpdateProduct)
	e.DELETE("/products/:id", handler.DeleteProduct)

	return &CQRSServiceServer{
		logger: logger,
		e:      e,
	}
}

// RegisterLifecycleHooks はサーバーのライフサイクルフックを登録します。
//
// Parameters:
//   - lc: fxライフサイクル
//   - server: CQRSServiceServer
//   - cfg: サーバー設定
func RegisterLifecycleHooks(lc fx.Lifecycle, server *CQRSServiceServer, cfg *CQRSServiceConfig) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// port競合によるエラーを回避し、動的ポート割り当てをサポートするため、
			// 事前にListenしてからStartServerを使用する
			ln, err := net.Listen("tcp", ":"+cfg.Port)
			if err != nil {
				return err
			}
			// 実際に割り当てられたアドレスを保存（ポート0の場合、動的に割り当てられる）
			server.Addr = ln.Addr().String()
			server.e.Listener = ln
			go func() {
				server.logger.InfoContext(ctx, "Starting CQRS Client Service Server", slog.String("addr", server.Addr))
				// NOTE: EchoのStartは内部でListenを呼び出すため、ここではStart("")を使用する
				// https://echo.labstack.com/docs/customization?utm_source=chatgpt.com#custom-listener
				if err = server.e.Start(""); err != nil && !errors.Is(err, http.ErrServerClosed) {
					server.logger.ErrorContext(ctx, "Failed to start server", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.logger.InfoContext(ctx, "Shutting down CQRS Client Service Server...")
			return server.e.Shutdown(ctx)
		},
	})
}
