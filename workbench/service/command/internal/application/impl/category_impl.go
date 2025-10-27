package impl

import (
	"context"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

// CategoryServiceImpl はカテゴリサービスの実装です。
//
// Fields:
//   - repo: カテゴリのデータ永続化を担当するリポジトリ
//   - tm: トランザクション管理を担当するマネージャー
//   - logger: 構造化ログ出力を担当するロガー
type CategoryServiceImpl struct {
	repo   categories.CategoryRepository
	tm     service.TransactionManager
	logger *slog.Logger
}

// NewCategoryServiceImpl は新しいCategoryServiceImplインスタンスを生成します。
//
// Parameters:
//   - logger: 構造化ログ出力用のロガー
//   - repo: カテゴリデータの永続化を担うリポジトリ
//   - tm: トランザクション管理を担うマネージャー
//
// Returns:
//   - *CategoryServiceImpl: 初期化されたカテゴリサービス実装
func NewCategoryServiceImpl(logger *slog.Logger, repo categories.CategoryRepository, tm service.TransactionManager) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		repo:   repo,
		tm:     tm,
		logger: logger,
	}
}

// Add は新しいカテゴリを追加します。
// カテゴリ名に重複がないことを確認してから、永続化します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - categoryDTO: 追加するカテゴリ情報
//
// Returns:
//   - *dto.CategoryDTO: 作成されたカテゴリのDTO
//   - error: カテゴリ名の重複や、その他の永続化に関するエラー
func (s *CategoryServiceImpl) Add(ctx context.Context, categoryDTO *dto.CreateCategoryDTO) (*dto.CategoryDTO, error) {
	category, err := dto.CategoryFromCreateDTO(categoryDTO)
	if err != nil {
		return nil, err
	}

	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// NOTE: defer内でerrを評価するため、クロージャで囲む。defer時点のerrを参照させるため。
	defer func() {
		if completeErr := s.tm.Complete(ctx, tx, err); completeErr != nil {
			s.logger.ErrorContext(ctx, "トランザクションの完了に失敗しました", slog.Any("error", completeErr))
		}
	}()

	exists, err := s.repo.ExistsByName(ctx, tx, category.Name())
	if err != nil {
		return nil, err
	}
	if exists {
		// defer内で評価されるerrに代入し、トランザクション完了時にロールバックされるようにする
		err = errs.NewApplicationError("CATEGORY_ALREADY_EXISTS", "Category already exists")
		return nil, err
	}

	if err = s.repo.Create(ctx, tx, category); err != nil {
		return nil, err
	}

	return dto.NewCategoryDTOFromEntity(category), nil
}

// Update は既存のカテゴリ情報を更新します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - categoryDTO: 更新するカテゴリ情報（ID、名前を含む）
//
// Returns:
//   - *dto.CategoryDTO: 更新されたカテゴリのDTO
//   - error: 指定したIDのカテゴリが存在しない場合や、その他の永続化に関するエラー
func (s *CategoryServiceImpl) Update(ctx context.Context, categoryDTO *dto.UpdateCategoryDTO) (*dto.CategoryDTO, error) {
	category, err := dto.CategoryFromUpdateDTO(categoryDTO)
	if err != nil {
		return nil, err
	}

	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if completeErr := s.tm.Complete(ctx, tx, err); completeErr != nil {
			s.logger.ErrorContext(ctx, "トランザクションの完了に失敗しました", slog.Any("error", completeErr))
		}
	}()

	if err = s.repo.UpdateById(ctx, tx, category); err != nil {
		return nil, err
	}

	return dto.NewCategoryDTOFromEntity(category), nil
}

// Delete は指定されたカテゴリを削除します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - categoryDTO: 削除対象カテゴリのID情報
//
// Returns:
//   - *dto.CategoryDTO: 削除されたカテゴリのDTO
//   - error: 指定したIDのカテゴリが存在しない場合や、その他の永続化に関するエラー
func (s *CategoryServiceImpl) Delete(ctx context.Context, categoryDTO *dto.DeleteCategoryDTO) (*dto.CategoryDTO, error) {
	categoryID, err := dto.CategoryIdFromDeleteDTO(categoryDTO)
	if err != nil {
		return nil, err
	}

	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if completeErr := s.tm.Complete(ctx, tx, err); completeErr != nil {
			s.logger.ErrorContext(ctx, "トランザクションの完了に失敗しました", slog.Any("error", completeErr))
		}
	}()

	category, err := s.repo.FindById(ctx, tx, categoryID)
	if err != nil {
		return nil, err
	}

	if err = s.repo.DeleteById(ctx, tx, categoryID); err != nil {
		return nil, err
	}

	return dto.NewCategoryDTOFromEntity(category), nil
}

var _ service.CategoryService = (*CategoryServiceImpl)(nil)
