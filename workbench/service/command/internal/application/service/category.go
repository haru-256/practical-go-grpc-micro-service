package service

import (
	"context"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
)

// CategoryService はカテゴリに関するアプリケーションサービスのインターフェースです。
// カテゴリの追加、更新、削除といったビジネスロジックを提供し、
// トランザクション管理とドメインルールの適用を担います。
type CategoryService interface {
	// Add は新しいカテゴリを追加します。
	// カテゴリ名の重複をチェックし、既に存在する場合はApplicationErrorを返します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - category: 追加するカテゴリ情報
	//
	// Returns:
	//   - error: カテゴリ名が既に存在する場合はApplicationError (コード: CATEGORY_ALREADY_EXISTS)、
	//     データベースエラーが発生した場合はCRUDError
	Add(ctx context.Context, category *categories.Category) error

	// Update は既存のカテゴリ情報を更新します。
	// リポジトリ層でカテゴリの存在確認を行い、存在しない場合はCRUDErrorを返します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - category: 更新するカテゴリ情報
	//
	// Returns:
	//   - error: カテゴリが存在しない場合はCRUDError (コード: NOT_FOUND)、
	//     データベースエラーが発生した場合はCRUDError
	Update(ctx context.Context, category *categories.Category) error

	// Delete は指定されたカテゴリを削除します。
	// リポジトリ層でカテゴリの存在確認を行い、存在しない場合はCRUDErrorを返します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - category: 削除するカテゴリ情報
	//
	// Returns:
	//   - error: カテゴリが存在しない場合はCRUDError (コード: NOT_FOUND)、
	//     データベースエラーが発生した場合はCRUDError
	Delete(ctx context.Context, category *categories.Category) error
}
