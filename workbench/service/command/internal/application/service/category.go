package service

import (
	"context"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
)

// CategoryService はカテゴリに関するアプリケーションサービスのインターフェースです。
//
//go:generate go tool mockgen -source=$GOFILE -destination=../../mock/service/category_service_mock.go -package=mock_service
type CategoryService interface {
	// Add は新しいカテゴリを追加します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - categoryDTO: 追加するカテゴリ情報
	//
	// Returns:
	//   - *dto.CategoryDTO: 作成されたカテゴリ
	//   - error: エラー
	Add(ctx context.Context, categoryDTO *dto.CreateCategoryDTO) (*dto.CategoryDTO, error)

	// Update は既存のカテゴリ情報を更新します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - categoryDTO: 更新するカテゴリ情報
	//
	// Returns:
	//   - *dto.CategoryDTO: 更新されたカテゴリ
	//   - error: エラー
	Update(ctx context.Context, categoryDTO *dto.UpdateCategoryDTO) (*dto.CategoryDTO, error)

	// Delete は指定されたカテゴリを削除します。
	//
	// Parameters:
	//   - ctx: コンテキスト
	//   - categoryDTO: 削除するカテゴリ情報
	//
	// Returns:
	//   - *dto.CategoryDTO: 削除されたカテゴリ
	//   - error: エラー
	Delete(ctx context.Context, categoryDTO *dto.DeleteCategoryDTO) (*dto.CategoryDTO, error)
}
