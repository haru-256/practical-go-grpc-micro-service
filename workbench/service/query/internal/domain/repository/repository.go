package repository

import (
	"context"

	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/domain/models"
)

type ProductRepository interface {
	// List はすべての商品を取得します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//
	// Returns:
	//   - []*models.Product: 商品リスト
	//   - error: エラー
	List(ctx context.Context) ([]*models.Product, error)

	// FindById は商品IDで商品を検索します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - id: 商品ID
	//
	// Returns:
	//   - *models.Product: 商品
	//   - error: エラー
	FindById(ctx context.Context, id string) (*models.Product, error)

	// FindByNameLike は商品名の部分一致で商品を検索します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - keyword: 検索キーワード
	//
	// Returns:
	//   - []*models.Product: 商品リスト
	//   - error: エラー
	FindByNameLike(ctx context.Context, keyword string) ([]*models.Product, error)
}

type CategoryRepository interface {
	// List はすべてのカテゴリを取得します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//
	// Returns:
	//   - []*models.Category: カテゴリリスト
	//   - error: エラー
	List(ctx context.Context) ([]*models.Category, error)

	// FindById はカテゴリIDでカテゴリを検索します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - id: カテゴリID
	//
	// Returns:
	//   - *models.Category: カテゴリ
	//   - error: エラー
	FindById(ctx context.Context, id string) (*models.Category, error)
}
