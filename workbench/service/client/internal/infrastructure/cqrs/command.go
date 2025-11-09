package cqrs

import (
	"net/http"

	"connectrpc.com/connect"
	cmdconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/command/v1/commandv1connect"
)

// CommandServiceClient はCommand Serviceへの接続を管理するクライアント
type CommandServiceClient struct {
	Category cmdconnect.CategoryServiceClient // カテゴリサービスクライアント
	Product  cmdconnect.ProductServiceClient  // 商品サービスクライアント
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

	return &CommandServiceClient{
		Category: categoryClient,
		Product:  productClient,
	}
}
