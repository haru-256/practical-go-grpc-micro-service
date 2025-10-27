package service

import (
	"context"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
)

// ProductService は商品に関するアプリケーションサービスのインターフェースです。
//
//go:generate go tool mockgen -source=$GOFILE -destination=./mock_product.go -package=service
type ProductService interface {
	// Add は新しい商品を追加します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - ProductDTO: 追加する商品情報
	//
	// Returns:
	//   - *dto.ProductDTO: 作成された商品
	//   - error: エラー
	Add(ctx context.Context, ProductDTO *dto.CreateProductDTO) (*dto.ProductDTO, error)

	// Update は既存の商品情報を更新します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - productDTO: 更新する商品情報
	//
	// Returns:
	//   - *dto.ProductDTO: 更新された商品
	//   - error: エラー
	Update(ctx context.Context, productDTO *dto.UpdateProductDTO) (*dto.ProductDTO, error)

	// Delete は指定された商品を削除します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - productDTO: 削除する商品情報
	//
	// Returns:
	//   - *dto.ProductDTO: 削除された商品
	//   - error: エラー
	Delete(ctx context.Context, productDTO *dto.DeleteProductDTO) (*dto.ProductDTO, error)
}
