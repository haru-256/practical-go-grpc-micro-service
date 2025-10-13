package products

import (
	"context"
	"database/sql"
)

// ProductRepository は商品エンティティの永続化を担うリポジトリインターフェースです。
// データベースへの商品データのCRUD操作を提供します。
type ProductRepository interface {
	// ExistsById は指定された商品IDが存在するかチェックします。
	ExistsById(ctx context.Context, tx *sql.Tx, id *ProductId) (bool, error)
	// ExistsByName は指定された商品名が存在するかチェックします。
	ExistsByName(ctx context.Context, tx *sql.Tx, name *ProductName) (bool, error)
	// Create は新しい商品を作成します。
	Create(ctx context.Context, tx *sql.Tx, product *Product) error
	// UpdateById は商品IDを指定して商品情報を更新します。
	UpdateById(ctx context.Context, tx *sql.Tx, product *ProductId) error
	// DeleteById は商品IDを指定して商品を削除します。
	DeleteById(ctx context.Context, tx *sql.Tx, id *ProductId) error
}
