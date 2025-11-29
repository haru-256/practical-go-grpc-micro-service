package cqrs

import (
	"context"
	"fmt"
	"net/http"

	"buf.build/gen/go/grpc/grpc/connectrpc/go/grpc/health/v1/healthv1connect"
	healthv1 "buf.build/gen/go/grpc/grpc/protocolbuffers/go/grpc/health/v1"
	"connectrpc.com/connect"
	cmdconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/command/v1/commandv1connect"
)

// CommandServiceClient はCommand Serviceへの接続を管理するクライアント
type CommandServiceClient struct {
	Category     cmdconnect.CategoryServiceClient // カテゴリサービスクライアント
	Product      cmdconnect.ProductServiceClient  // 商品サービスクライアント
	healthClient healthv1connect.HealthClient     // ヘルスチェッククライアント
	serviceURL   string                           // サービスURL
}

// NewCommandServiceClient はCommandServiceClientを生成します。
//
// Parameters:
//   - client: HTTPクライアント
//   - cfg: CQRS設定
//
// Returns:
//   - *CommandServiceClient: CommandServiceClient
func NewCommandServiceClient(client *http.Client, cfg *CQRSServiceConfig) *CommandServiceClient {
	categoryClient := cmdconnect.NewCategoryServiceClient(client, cfg.CommandServiceURL, connect.WithGRPC())
	productClient := cmdconnect.NewProductServiceClient(client, cfg.CommandServiceURL, connect.WithGRPC())
	healthClient := healthv1connect.NewHealthClient(client, cfg.CommandServiceURL, connect.WithGRPC())

	return &CommandServiceClient{
		Category:     categoryClient,
		Product:      productClient,
		healthClient: healthClient,
		serviceURL:   cfg.CommandServiceURL,
	}
}

// HealthCheck はCommand Serviceへのヘルスチェックを実行します。
//
// Parameters:
//   - ctx: コンテキスト
//
// Returns:
//   - error: ヘルスチェックエラー
func (c *CommandServiceClient) HealthCheck(ctx context.Context) error {
	req := connect.NewRequest(&healthv1.HealthCheckRequest{})
	resp, err := c.healthClient.Check(ctx, req)
	if err != nil {
		return fmt.Errorf("health check failed for %s: %w", c.serviceURL, err)
	}

	if resp.Msg.Status != healthv1.HealthCheckResponse_SERVING {
		return fmt.Errorf("service not healthy: %s (status: %s)", c.serviceURL, resp.Msg.Status.String())
	}

	return nil
}
