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
	// 同じ名前のカテゴリが既に存在する場合はApplicationErrorを返します。
	Add(ctx context.Context, category *categories.Category) error

	// Update は既存のカテゴリ情報を更新します。
	// 指定されたIDのカテゴリが存在しない場合はCRUDErrorを返します。
	Update(ctx context.Context, category *categories.Category) error

	// Delete は指定されたカテゴリを削除します。
	// 指定されたIDのカテゴリが存在しない場合はCRUDErrorを返します。
	Delete(ctx context.Context, category *categories.Category) error
}
