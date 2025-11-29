package cqrs

import (
	"context"
	"fmt"
	"net/http"

	"buf.build/gen/go/grpc/grpc/connectrpc/go/grpc/health/v1/healthv1connect"
	healthv1 "buf.build/gen/go/grpc/grpc/protocolbuffers/go/grpc/health/v1"
	"connectrpc.com/connect"
	queryconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1/queryv1connect"
)

// QueryServiceClient はQuery Serviceへの接続を管理するクライアント
type QueryServiceClient struct {
	Category     queryconnect.CategoryServiceClient // カテゴリサービスクライアント
	Product      queryconnect.ProductServiceClient  // 商品サービスクライアント
	healthClient healthv1connect.HealthClient       // ヘルスチェッククライアント
	serviceURL   string                             // サービスURL
}

// NewQueryServiceClient はQueryServiceClientを生成します。
//
// Parameters:
//   - client: HTTPクライアント
//   - cfg: CQRS設定
//
// Returns:
//   - *QueryServiceClient: QueryServiceClient
func NewQueryServiceClient(client *http.Client, cfg *CQRSServiceConfig) *QueryServiceClient {
	categoryClient := queryconnect.NewCategoryServiceClient(client, cfg.QueryServiceURL, connect.WithGRPC())
	productClient := queryconnect.NewProductServiceClient(client, cfg.QueryServiceURL, connect.WithGRPC())
	healthClient := healthv1connect.NewHealthClient(client, cfg.QueryServiceURL, connect.WithGRPC())

	return &QueryServiceClient{
		Category:     categoryClient,
		Product:      productClient,
		healthClient: healthClient,
		serviceURL:   cfg.QueryServiceURL,
	}
}

// HealthCheck はQuery Serviceへのヘルスチェックを実行します。
//
// Parameters:
//   - ctx: コンテキスト
//
// Returns:
//   - error: ヘルスチェックエラー
func (c *QueryServiceClient) HealthCheck(ctx context.Context) error {
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
