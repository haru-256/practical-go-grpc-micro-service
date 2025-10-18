package products

import (
	"context"
	"database/sql"
)

// ProductRepository は商品エンティティの永続化を担うリポジトリインターフェースです。
// データベースへの商品データのCRUD操作を提供します。
//
//go:generate go tool mockgen -source=$GOFILE -destination=./mock_product_repository.go -package=products
type ProductRepository interface {
	// ExistsById は指定された商品IDが存在するかチェックします。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - id: チェック対象の商品ID
	//
	// Returns:
	//   - bool: 商品が存在する場合はtrue、存在しない場合はfalse
	//   - error: データベースエラーが発生した場合
	ExistsById(ctx context.Context, tx *sql.Tx, id *ProductId) (bool, error)

	// ExistsByName は指定された商品名が存在するかチェックします。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - name: チェック対象の商品名
	//
	// Returns:
	//   - bool: 商品が存在する場合はtrue、存在しない場合はfalse
	//   - error: データベースエラーが発生した場合
	ExistsByName(ctx context.Context, tx *sql.Tx, name *ProductName) (bool, error)

	// Create は新しい商品をデータベースに追加します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - product: 追加する商品エンティティ
	//
	// Returns:
	//   - error: 追加に失敗した場合のエラー
	Create(ctx context.Context, tx *sql.Tx, product *Product) error

	// UpdateById は商品IDを指定して商品情報を更新します。
	// 指定されたIDの商品が存在しない場合はCRUDError (NOT_FOUND)を返します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - product: 更新する商品エンティティ（IDで識別）
	//
	// Returns:
	//   - error: 更新に失敗した場合のエラー
	//     - CRUDError (NOT_FOUND): 商品が見つからない場合
	UpdateById(ctx context.Context, tx *sql.Tx, product *Product) error

	// DeleteById は商品IDを指定して商品を削除します。
	// 指定されたIDの商品が存在しない場合はCRUDError (NOT_FOUND)を返します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - id: 削除する商品のID
	//
	// Returns:
	//   - error: 削除に失敗した場合のエラー
	//     - CRUDError (NOT_FOUND): 商品が見つからない場合
	DeleteById(ctx context.Context, tx *sql.Tx, id *ProductId) error
}
