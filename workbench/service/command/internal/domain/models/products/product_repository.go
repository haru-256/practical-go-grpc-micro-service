package products

import (
	"context"
	"database/sql"
)

// ProductRepository は商品エンティティの永続化を担うリポジトリインターフェースです。
//
//go:generate go tool mockgen -source=$GOFILE -destination=../../../mock/repository/product_repository_mock.go -package=mock_repository
type ProductRepository interface {
	// ExistsById は指定された商品IDが存在するかチェックします。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - id: 商品ID
	//
	// Returns:
	//   - bool: 商品が存在する場合はtrue
	//   - error: エラー
	ExistsById(ctx context.Context, tx *sql.Tx, id *ProductId) (bool, error)

	// FindById は指定された商品IDで商品を取得します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - id: 商品ID
	//
	// Returns:
	//   - *Product: 商品エンティティ
	//   - error: エラー
	FindById(ctx context.Context, tx *sql.Tx, id *ProductId) (*Product, error)

	// ExistsByName は指定された商品名が存在するかチェックします。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - name: 商品名
	//
	// Returns:
	//   - bool: 商品が存在する場合はtrue
	//   - error: エラー
	ExistsByName(ctx context.Context, tx *sql.Tx, name *ProductName) (bool, error)

	// Create は新しい商品をデータベースに追加します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - product: 商品エンティティ
	//
	// Returns:
	//   - error: エラー
	Create(ctx context.Context, tx *sql.Tx, product *Product) error

	// UpdateById は商品IDを指定して商品情報を更新します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - product: 商品エンティティ
	//
	// Returns:
	//   - error: エラー
	UpdateById(ctx context.Context, tx *sql.Tx, product *Product) error

	// DeleteById は商品IDを指定して商品を削除します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - id: 商品ID
	//
	// Returns:
	//   - error: エラー
	DeleteById(ctx context.Context, tx *sql.Tx, id *ProductId) error
}
