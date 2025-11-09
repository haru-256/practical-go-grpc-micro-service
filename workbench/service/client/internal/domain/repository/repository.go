package repository

import (
	"context"

	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/domain/models"
)

// CQRSRepository はCQRSパターンに基づくリポジトリインターフェース
// Command ServiceとQuery Serviceへの書き込み・読み取り操作を提供します。
//
//go:generate go tool mockgen -source=$GOFILE -destination=../../mock/repository/repository_mock.go -package=mock_repository
type CQRSRepository interface {
	// CreateCategory はカテゴリを作成します。
	CreateCategory(ctx context.Context, categoryName string) (*models.Category, error)
	// UpdateCategory はカテゴリを更新します。
	UpdateCategory(ctx context.Context, category *models.Category) (*models.Category, error)
	// DeleteCategory はカテゴリを削除します。
	DeleteCategory(ctx context.Context, id string) error
	// CategoryList はカテゴリ一覧を取得します。
	CategoryList(ctx context.Context) ([]*models.Category, error)
	// CategoryById はIDでカテゴリを取得します。
	CategoryById(ctx context.Context, id string) (*models.Category, error)

	// CreateProduct は商品を作成します。
	CreateProduct(ctx context.Context, productName string, productPrice uint32, category *models.Category) (*models.Product, error)
	// UpdateProduct は商品を更新します。
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	// DeleteProduct は商品を削除します。
	DeleteProduct(ctx context.Context, id string) error
	// ProductList は商品一覧を取得します。
	ProductList(ctx context.Context) ([]*models.Product, error)
	// ProductById はIDで商品を取得します。
	ProductById(ctx context.Context, id string) (*models.Product, error)
	// ProductByKeyword はキーワードで商品を検索します。
	ProductByKeyword(ctx context.Context, keyword string) ([]*models.Product, error)
}
