package cqrs

import (
	"net/http"

	"connectrpc.com/connect"
	queryconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1/queryv1connect"
)

// QueryServiceClient はQuery Serviceへの接続を管理するクライアント
type QueryServiceClient struct {
	Category queryconnect.CategoryServiceClient // カテゴリサービスクライアント
	Product  queryconnect.ProductServiceClient  // 商品サービスクライアント
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

	return &QueryServiceClient{
		Category: categoryClient,
		Product:  productClient,
	}
}
