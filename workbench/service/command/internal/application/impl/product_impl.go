package impl

import (
	"context"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/errs"
)

// ProductServiceImpl は商品サービスの実装です。
type ProductServiceImpl struct {
	repo   products.ProductRepository
	tm     service.TransactionManager
	logger *slog.Logger
}

// NewProductServiceImpl は新しいProductServiceImplインスタンスを生成します。
// インターフェースを受け入れ、具象型を返すことで、柔軟性と明確性を両立させます。
//
// 使用例:
//
//	svc := impl.NewProductServiceImpl(repo, tm)
//	var productService service.ProductService = svc  // インターフェースとして使用
func NewProductServiceImpl(logger *slog.Logger, repo products.ProductRepository, tm service.TransactionManager) *ProductServiceImpl {
	return &ProductServiceImpl{
		repo:   repo,
		tm:     tm,
		logger: logger,
	}
}

// Add は新しい商品を追加します。
// 商品名の重複をチェックし、既に存在する場合はApplicationErrorを返します。
// トランザクション内で実行され、エラー時は自動的にロールバックされます。
//
// Parameters:
//   - ctx: コンテキスト
//   - product: 追加する商品エンティティ
//
// Returns:
//   - error: 追加に失敗した場合のエラー
func (s *ProductServiceImpl) Add(ctx context.Context, product *products.Product) error {
	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return err
	}
	// NOTE: defer内でerrを評価するため、クロージャで囲む。defer時点のerrを参照させるため。
	defer func() {
		if completeErr := s.tm.Complete(ctx, tx, err); completeErr != nil {
			s.logger.ErrorContext(ctx, "トランザクションの完了に失敗しました", slog.Any("error", completeErr))
		}
	}()

	exists, err := s.repo.ExistsByName(ctx, tx, product.Name())
	if err != nil {
		return err
	}
	if exists {
		// defer内で評価されるerrに代入し、トランザクション完了時にロールバックされるようにする
		err = errs.NewApplicationError("PRODUCT_ALREADY_EXISTS", "Product already exists")
		return err
	}

	if err = s.repo.Create(ctx, tx, product); err != nil {
		return err
	}

	return nil
}

// Update は既存の商品情報を更新します。
// リポジトリ層で商品の存在確認を行います。
// トランザクション内で実行され、エラー時は自動的にロールバックされます。
//
// Parameters:
//   - ctx: コンテキスト
//   - product: 更新する商品エンティティ（IDで識別）
//
// Returns:
//   - error: 更新に失敗した場合のエラー（リポジトリ層から伝播）
func (s *ProductServiceImpl) Update(ctx context.Context, product *products.Product) error {
	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if completeErr := s.tm.Complete(ctx, tx, err); completeErr != nil {
			s.logger.ErrorContext(ctx, "トランザクションの完了に失敗しました", slog.Any("error", completeErr))
		}
	}()

	if err = s.repo.UpdateById(ctx, tx, product); err != nil {
		return err
	}

	return nil
}

// Delete は指定された商品を削除します。
// リポジトリ層で商品の存在確認を行います。
// トランザクション内で実行され、エラー時は自動的にロールバックされます。
//
// Parameters:
//   - ctx: コンテキスト
//   - product: 削除する商品エンティティ（IDで識別）
//
// Returns:
//   - error: 削除に失敗した場合のエラー（リポジトリ層から伝播）
func (s *ProductServiceImpl) Delete(ctx context.Context, product *products.Product) error {
	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if completeErr := s.tm.Complete(ctx, tx, err); completeErr != nil {
			s.logger.ErrorContext(ctx, "トランザクションの完了に失敗しました", slog.Any("error", completeErr))
		}
	}()

	if err = s.repo.DeleteById(ctx, tx, product.Id()); err != nil {
		return err
	}

	return nil
}
