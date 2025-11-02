package impl

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
)

// ProductServiceImpl は商品サービスの実装です。
// 商品のCRUD操作とトランザクション管理を提供します。
type ProductServiceImpl struct {
	productRepo  products.ProductRepository    // 商品データの永続化
	categoryRepo categories.CategoryRepository // カテゴリデータの取得
	tm           service.TransactionManager    // トランザクション管理
	logger       *slog.Logger                  // 構造化ログ出力
}

// NewProductServiceImpl は新しいProductServiceImplインスタンスを生成します。
//
// Parameters:
//   - logger: 構造化ログ出力用のロガー
//   - productRepo: 商品データの永続化を担うリポジトリ
//   - categoryRepo: カテゴリデータの取得を担うリポジトリ
//   - tm: トランザクション管理を担うマネージャー
//
// Returns:
//   - *ProductServiceImpl: 初期化された商品サービス実装
func NewProductServiceImpl(logger *slog.Logger, productRepo products.ProductRepository, categoryRepo categories.CategoryRepository, tm service.TransactionManager) *ProductServiceImpl {
	return &ProductServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		tm:           tm,
		logger:       logger,
	}
}

// Add は新しい商品を追加します。
// 商品名の重複チェックを行い、トランザクション内で作成処理を実行します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - productDTO: 追加する商品情報
//
// Returns:
//   - *dto.ProductDTO: 作成された商品のDTO
//   - error: エラー情報
func (s *ProductServiceImpl) Add(ctx context.Context, productDTO *dto.CreateProductDTO) (result *dto.ProductDTO, err error) {
	var (
		tx      *sql.Tx
		product *products.Product
		exists  bool
	)

	product, err = dto.ProductFromCreateDTO(productDTO)
	if err != nil {
		return nil, err
	}

	tx, err = s.tm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// NOTE: defer内でerrを評価するため、クロージャで囲む。defer時点のerrを参照させるため。
	defer func() {
		handleTransactionComplete(ctx, s.tm, tx, &err, &result, s.logger)
	}()

	// categoryが既に存在するかチェック

	exists, err = s.productRepo.ExistsByName(ctx, tx, product.Name())
	if err != nil {
		return nil, err
	}
	if exists {
		// defer内で評価されるerrに代入し、トランザクション完了時にロールバックされるようにする
		err = errs.NewApplicationError("PRODUCT_ALREADY_EXISTS", "Product already exists")
		return nil, err
	}

	if err = s.productRepo.Create(ctx, tx, product); err != nil {
		return nil, err
	}

	result = dto.NewProductDTOFromEntity(product)
	return result, nil
}

// Update は既存の商品情報を更新します。
// カテゴリ情報を取得してからトランザクション内で更新処理を実行します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - productDTO: 更新する商品情報
//
// Returns:
//   - *dto.ProductDTO: 更新された商品のDTO
//   - error: エラー情報
func (s *ProductServiceImpl) Update(ctx context.Context, productDTO *dto.UpdateProductDTO) (result *dto.ProductDTO, err error) {
	var (
		tx         *sql.Tx
		categoryId *categories.CategoryId
		category   *categories.Category
		product    *products.Product
	)

	tx, err = s.tm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		handleTransactionComplete(ctx, s.tm, tx, &err, &result, s.logger)
	}()

	// カテゴリ名をrepositoryから取得
	// FIXME: カテゴリ名は更新しないのに毎回取得するのは無駄がある
	categoryId, err = categories.NewCategoryId(productDTO.CategoryId)
	if err != nil {
		return nil, err
	}
	category, err = s.categoryRepo.FindById(ctx, tx, categoryId)
	if err != nil {
		return nil, err
	}

	product, err = dto.ProductFromUpdateDTO(productDTO, category.Name())
	if err != nil {
		return nil, err
	}

	if err = s.productRepo.UpdateById(ctx, tx, product); err != nil {
		return nil, err
	}

	result = dto.NewProductDTOFromEntity(product)
	return result, nil
}

// Delete は指定された商品を削除します。
// 削除前に商品の存在確認を行い、トランザクション内で削除処理を実行します。
//
// Parameters:
//   - ctx: リクエストコンテキスト
//   - productDTO: 削除対象商品のID情報を含むDTO
//
// Returns:
//   - *dto.ProductDTO: 削除された商品のDTO
//   - error: エラー情報
func (s *ProductServiceImpl) Delete(ctx context.Context, productDTO *dto.DeleteProductDTO) (result *dto.ProductDTO, err error) {
	var (
		product   *products.Product
		productID *products.ProductId
	)
	productID, err = dto.ProductIdFromDeleteDTO(productDTO)
	if err != nil {
		return nil, err
	}

	tx, err := s.tm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		handleTransactionComplete(ctx, s.tm, tx, &err, &result, s.logger)
	}()

	product, err = s.productRepo.FindById(ctx, tx, productID)
	if err != nil {
		return nil, err
	}

	if err = s.productRepo.DeleteById(ctx, tx, productID); err != nil {
		return nil, err
	}

	result = dto.NewProductDTOFromEntity(product)
	return result, nil
}

var _ service.ProductService = (*ProductServiceImpl)(nil)
