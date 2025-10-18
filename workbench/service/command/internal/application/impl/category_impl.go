package impl

import (
	"context"
	"log"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

// CategoryServiceImpl はカテゴリサービスの実装です。
type CategoryServiceImpl struct {
	repo categories.CategoryRepository
	tm   service.TransactionManager
}

// NewCategoryServiceImpl は新しいCategoryServiceImplインスタンスを生成します。
// インターフェースを受け入れ、具象型を返すことで、柔軟性と明確性を両立させます。
//
// 使用例:
//
//	svc := impl.NewCategoryServiceImpl(repo, tm)
//	var categoryService service.CategoryService = svc  // インターフェースとして使用
func NewCategoryServiceImpl(repo categories.CategoryRepository, tm service.TransactionManager) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		repo: repo,
		tm:   tm,
	}
}

// Add は新しいカテゴリを追加します。
// カテゴリ名の重複をチェックし、既に存在する場合はApplicationErrorを返します。
// トランザクション管理を行い、エラー時は自動的にロールバックされます。
//
// Parameters:
//   - ctx: コンテキスト
//   - category: 追加するカテゴリ情報
//
// Returns:
//   - error: カテゴリ名が既に存在する場合はApplicationError (コード: CATEGORY_ALREADY_EXISTS)、
//     データベースエラーが発生した場合はCRUDError
func (s *CategoryServiceImpl) Add(ctx context.Context, category *categories.Category) error {
	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return err
	}
	// NOTE: defer内でerrを評価するため、クロージャで囲む。defer時点のerrを参照させるため。
	defer func() {
		if completeErr := s.tm.Complete(tx, err); completeErr != nil {
			log.Println("トランザクションの完了に失敗しました:", completeErr)
		}
	}()

	exists, err := s.repo.ExistsByName(ctx, tx, category.Name())
	if err != nil {
		return err
	}
	if exists {
		// defer内で評価されるerrに代入し、トランザクション完了時にロールバックされるようにする
		err = errs.NewApplicationError("CATEGORY_ALREADY_EXISTS", "Category already exists")
		return err
	}

	if err = s.repo.Create(ctx, tx, category); err != nil {
		return err
	}

	return nil
}

// Update は既存のカテゴリ情報を更新します。
// リポジトリ層でカテゴリの存在確認を行い、存在しない場合はCRUDErrorを返します。
// トランザクション管理を行い、エラー時は自動的にロールバックされます。
//
// Parameters:
//   - ctx: コンテキスト
//   - category: 更新するカテゴリ情報
//
// Returns:
//   - error: カテゴリが存在しない場合はCRUDError (コード: NOT_FOUND)、
//     データベースエラーが発生した場合はCRUDError
func (s *CategoryServiceImpl) Update(ctx context.Context, category *categories.Category) error {
	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if completeErr := s.tm.Complete(tx, err); completeErr != nil {
			log.Println("トランザクションの完了に失敗しました:", completeErr)
		}
	}()

	if err = s.repo.UpdateById(ctx, tx, category); err != nil {
		return err
	}

	return nil
}

// Delete は指定されたカテゴリを削除します。
// リポジトリ層でカテゴリの存在確認を行い、存在しない場合はCRUDErrorを返します。
// トランザクション管理を行い、エラー時は自動的にロールバックされます。
//
// Parameters:
//   - ctx: コンテキスト
//   - category: 削除するカテゴリ情報
//
// Returns:
//   - error: カテゴリが存在しない場合はCRUDError (コード: NOT_FOUND)、
//     データベースエラーが発生した場合はCRUDError
func (s *CategoryServiceImpl) Delete(ctx context.Context, category *categories.Category) error {
	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = s.tm.Complete(tx, err)
		if err != nil {
			log.Println("トランザクションの完了に失敗しました:", err)
		}
	}()

	if err = s.repo.DeleteById(ctx, tx, category.Id()); err != nil {
		return err
	}

	return nil
}
