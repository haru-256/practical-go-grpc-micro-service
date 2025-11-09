package cqrs

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/utils"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// CQRSServiceConfig はCQRSサービスの接続設定
type CQRSServiceConfig struct {
	CommandServiceURL     string        // CommandサービスのURL
	QueryServiceURL       string        // QueryサービスのURL
	RequestTimeout        time.Duration // リクエストタイムアウト
	TCPTimeout            time.Duration // TCP接続タイムアウト
	TCPKeepAlive          time.Duration // TCPキープアライブ間隔
	TLSHandshakeTimeout   time.Duration // TLSハンドシェイクタイムアウト
	ResponseHeaderTimeout time.Duration // レスポンスヘッダー受信タイムアウト
	IdleConnTimeout       time.Duration // アイドル接続タイムアウト
	MaxIdleConns          int           // 最大アイドル接続数
	MaxIdleConnsPerHost   int           // ホストごとの最大アイドル接続数
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
		CommandServiceURL:     utils.GetKey[string](v, "cqrs.command_service_url", &configErrors),
		QueryServiceURL:       utils.GetKey[string](v, "cqrs.query_service_url", &configErrors),
		RequestTimeout:        utils.GetKey[time.Duration](v, "cqrs.request_timeout", &configErrors),
		TCPTimeout:            utils.GetKey[time.Duration](v, "cqrs.tcp_timeout", &configErrors),
		TCPKeepAlive:          utils.GetKey[time.Duration](v, "cqrs.tcp_keep_alive", &configErrors),
		TLSHandshakeTimeout:   utils.GetKey[time.Duration](v, "cqrs.tls_handshake_timeout", &configErrors),
		ResponseHeaderTimeout: utils.GetKey[time.Duration](v, "cqrs.response_header_timeout", &configErrors),
		IdleConnTimeout:       utils.GetKey[time.Duration](v, "cqrs.idle_conn_timeout", &configErrors),
		MaxIdleConns:          utils.GetKey[int](v, "cqrs.max_idle_conns", &configErrors),
		MaxIdleConnsPerHost:   utils.GetKey[int](v, "cqrs.max_idle_conns_per_host", &configErrors),
	}
	// すべての環境変数を読み込んだ後、エラーがあればまとめて返す
	if len(configErrors) > 0 {
		return cfg, errors.Join(configErrors...)
	}
	return cfg, nil
}

// NewClient はHTTPクライアントを生成します。
//
// Parameters:
//   - cfg: CQRS設定
//
// Returns:
//   - *http.Client: HTTPクライアント
func NewClient(cfg *CQRSServiceConfig) *http.Client {
	// 1. Transport の設定
	// (http.DefaultTransport をコピーしてカスタマイズするのが一般的)
	tr := http.DefaultTransport.(*http.Transport).Clone()

	// 2. 接続タイムアウト (net.Dialer)
	tr.DialContext = (&net.Dialer{
		Timeout:   cfg.TCPTimeout,
		KeepAlive: cfg.TCPKeepAlive,
	}).DialContext

	// 3. 詳細タイムアウト
	tr.TLSHandshakeTimeout = cfg.TLSHandshakeTimeout
	tr.ResponseHeaderTimeout = cfg.ResponseHeaderTimeout
	tr.IdleConnTimeout = cfg.IdleConnTimeout

	// 4. 接続プーリング
	tr.MaxIdleConns = cfg.MaxIdleConns               // 全体でのアイドル接続数
	tr.MaxIdleConnsPerHost = cfg.MaxIdleConnsPerHost // ホストごとのアイドル接続数

	// (オプション: MaxConnsPerHost で同時接続数を制限することも可能)
	// tr.MaxConnsPerHost = 10

	// 5. Client の設定
	client := &http.Client{
		// リクエスト全体のタイムアウト
		Timeout: cfg.RequestTimeout,

		// カスタマイズした Transport を設定
		Transport: tr,
	}

	return client
}

// RegisterLifecycleHooks はHTTPクライアントのライフサイクルフックを登録します。
//
// Parameters:
//   - lc: fxライフサイクル
//   - client: HTTPクライアント
//   - logger: ロガー
func RegisterLifecycleHooks(lc fx.Lifecycle, client *http.Client, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// アイドル接続をクローズ
			if t, ok := client.Transport.(*http.Transport); ok {
				logger.InfoContext(ctx, "Closing HTTP client idle connections...")
				t.CloseIdleConnections()
				logger.InfoContext(ctx, "HTTP client idle connections closed")
			} else {
				logger.WarnContext(ctx, "HTTP client transport is not of type *http.Transport; cannot close idle connections")
			}
			return nil
		},
	})
}
