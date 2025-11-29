package categories

import (
	"context"
	"database/sql"
)

// CategoryRepository はカテゴリエンティティの永続化を担うリポジトリインターフェースです。
// データベースへのカテゴリデータのCRUD操作を提供します。
//
//go:generate go tool mockgen -source=$GOFILE -destination=../../../mock/repository/category_repository_mock.go -package=mock_repository
type CategoryRepository interface {
	// ExistsByName は指定されたカテゴリ名が既に存在するかをチェックします。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - name: チェックするカテゴリ名
	//
	// Returns:
	//   - bool: カテゴリが存在する場合はtrue、存在しない場合はfalse
	//   - error: データベースエラーが発生した場合
	ExistsByName(ctx context.Context, tx *sql.Tx, name *CategoryName) (bool, error)

	// FindById は指定されたIDのカテゴリを取得します。
	// カテゴリが存在しない場合はNOT_FOUNDエラーを返します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - id: 取得するカテゴリのID
	//
	// Returns:
	//   - *Category: 取得したカテゴリエンティティ
	//   - error: カテゴリが存在しない場合はCRUDError (コード: NOT_FOUND)、
	//     データベースエラーが発生した場合はそのエラー
	FindById(ctx context.Context, tx *sql.Tx, id *CategoryId) (*Category, error)

	// Create は新しいカテゴリをデータベースに作成します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - category: 作成するカテゴリ情報
	//
	// Returns:
	//   - error: データベースエラーが発生した場合
	Create(ctx context.Context, tx *sql.Tx, category *Category) error

	// UpdateById は指定されたIDのカテゴリ情報を更新します。
	// カテゴリが存在しない場合はNOT_FOUNDエラーを返します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - category: 更新するカテゴリ情報
	//
	// Returns:
	//   - error: カテゴリが存在しない場合はCRUDError (コード: NOT_FOUND)、
	//     データベースエラーが発生した場合はそのエラー
	UpdateById(ctx context.Context, tx *sql.Tx, category *Category) error

	// DeleteById は指定されたIDのカテゴリを削除します。
	// カテゴリが存在しない場合はNOT_FOUNDエラーを返します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - id: 削除するカテゴリのID
	//
	// Returns:
	//   - error: カテゴリが存在しない場合はCRUDError (コード: NOT_FOUND)、
	//     データベースエラーが発生した場合はそのエラー
	DeleteById(ctx context.Context, tx *sql.Tx, id *CategoryId) error

	// DeleteByName は指定されたカテゴリ名のカテゴリを削除します。
	// カテゴリが存在しない場合はNOT_FOUNDエラーを返します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - tx: トランザクション
	//   - name: 削除するカテゴリ名
	//
	// Returns:
	//   - error: カテゴリが存在しない場合はCRUDError (コード: NOT_FOUND)、
	//     データベースエラーが発生した場合はそのエラー
	DeleteByName(ctx context.Context, tx *sql.Tx, name *CategoryName) error
}
