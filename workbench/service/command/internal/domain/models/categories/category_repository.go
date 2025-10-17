package categories

import (
	"context"
	"database/sql"
)

// CategoryRepository はカテゴリエンティティの永続化を担うリポジトリインターフェースです。
// データベースへのカテゴリデータのCRUD操作を提供します。
type CategoryRepository interface {
	// ExistsByName は指定されたカテゴリ名が存在するかチェックします。
	ExistsByName(ctx context.Context, tx *sql.Tx, name *CategoryName) (bool, error)
	// Create は新しいカテゴリを作成します。
	Create(ctx context.Context, tx *sql.Tx, category *Category) error
	// UpdateById はカテゴリIDを指定してカテゴリ情報を更新します。
	UpdateById(ctx context.Context, tx *sql.Tx, category *Category) error
	// DeleteById はカテゴリIDを指定してカテゴリを削除します。
	DeleteById(ctx context.Context, tx *sql.Tx, id *CategoryId) error
	// DeleteByName はカテゴリ名を指定してカテゴリを削除します。
	DeleteByName(ctx context.Context, tx *sql.Tx, name *CategoryName) error
}
