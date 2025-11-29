package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/handler"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/models"
)

// ProductRepositoryImpl は商品リポジトリのSQLBoilerを使用した実装です。
type ProductRepositoryImpl struct {
	logger *slog.Logger
}

// NewProductRepositoryImpl は新しいProductRepositoryImplインスタンスを生成します。
func NewProductRepositoryImpl(logger *slog.Logger) *ProductRepositoryImpl {
	models.AddProductHook(boil.AfterInsertHook, productHookFactory(boil.AfterInsertHook, logger))
	models.AddProductHook(boil.AfterUpdateHook, productHookFactory(boil.AfterUpdateHook, logger))
	models.AddProductHook(boil.AfterDeleteHook, productHookFactory(boil.AfterDeleteHook, logger))
	return &ProductRepositoryImpl{logger: logger}
}

// ExistsById は指定された商品IDが存在するかどうかをチェックします。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - id: 商品ID
//
// Returns:
//   - bool: 商品が存在する場合はtrue
//   - error: データベースエラー
func (r *ProductRepositoryImpl) ExistsById(ctx context.Context, tx *sql.Tx, id *products.ProductId) (bool, error) {
	condition := models.ProductWhere.ObjID.EQ(id.Value())
	exists, err := models.Products(condition).Exists(ctx, tx)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to check if product exists", slog.Any("error", err))
		return false, handler.DBErrHandler(err)
	}
	return exists, nil
}

// FindById は指定されたIDの商品を取得します。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - id: 商品ID
//
// Returns:
//   - *products.Product: 見つかった商品エンティティ
//   - error: 商品が存在しない場合やデータベースエラー
func (r *ProductRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id *products.ProductId) (*products.Product, error) {
	condition := models.ProductWhere.ObjID.EQ(id.Value())
	model, err := models.Products(condition, qm.Load(models.ProductRels.Category)).One(ctx, tx)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to find product by ID", slog.Any("error", err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewCRUDError("NOT_FOUND", fmt.Sprintf("商品番号: %s は存在しないため、削除できませんでした。", id.Value()))
		} else {
			return nil, handler.DBErrHandler(err)
		}
	}
	// model.R.Category でCategoryにアクセスできる
	if model.R == nil || model.R.Category == nil {
		return nil, errs.NewCRUDError("NOT_FOUND", "関連するカテゴリが見つかりません")
	}
	// Categoryエンティティを構築
	categoryId, err := categories.NewCategoryId(model.R.Category.ObjID)
	if err != nil {
		return nil, err
	}
	categoryName, err := categories.NewCategoryName(model.R.Category.Name)
	if err != nil {
		return nil, err
	}
	category, err := categories.BuildCategory(categoryId, categoryName)
	if err != nil {
		return nil, err
	}

	// Productエンティティを構築
	productId, err := products.NewProductId(model.ObjID)
	if err != nil {
		return nil, err
	}
	productName, err := products.NewProductName(model.Name)
	if err != nil {
		return nil, err
	}
	productPrice, err := products.NewProductPrice(uint32(model.Price))
	if err != nil {
		return nil, err
	}

	product, err := products.BuildProduct(productId, productName, productPrice, category)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// ExistsByName は指定された名前の商品が存在するかどうかをチェックします。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - name: 商品名
//
// Returns:
//   - bool: 商品が存在する場合はtrue
//   - error: データベースエラー
func (r *ProductRepositoryImpl) ExistsByName(ctx context.Context, tx *sql.Tx, name *products.ProductName) (bool, error) {
	condition := models.ProductWhere.Name.EQ(name.Value())
	exists, err := models.Products(condition).Exists(ctx, tx)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to check if product exists", slog.Any("error", err))
		return false, handler.DBErrHandler(err)
	}
	return exists, nil
}

// Create は新しい商品をデータベースに追加します。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - product: 作成する商品
//
// Returns:
//   - error: データベースエラー
func (r *ProductRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, product *products.Product) error {
	newProduct := models.Product{
		ObjID:      product.Id().Value(),
		Name:       product.Name().Value(),
		Price:      int(product.Price().Value()),
		CategoryID: product.Category().Id().Value(),
	}
	// NOTE: boil.Infer() でauto-incrementのIDは無視され、勝手にDB側で採番された後、sqlboiler側の構造体にセットされる
	if err := newProduct.Insert(ctx, tx, boil.Infer()); err != nil {
		r.logger.ErrorContext(ctx, "Failed to create product", slog.Any("error", err))
		return handler.DBErrHandler(err)
	}

	return nil
}

// UpdateById は指定されたIDの商品情報を更新します。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - product: 更新する商品情報
//
// Returns:
//   - error: 商品が存在しない場合やデータベースエラー
func (r *ProductRepositoryImpl) UpdateById(ctx context.Context, tx *sql.Tx, Product *products.Product) error {
	condition := models.ProductWhere.ObjID.EQ(Product.Id().Value())
	upModel, err := models.Products(condition).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewCRUDErrorWithCause("NOT_FOUND", fmt.Sprintf("商品番号: %s は存在しないため、更新できませんでした。", Product.Id().Value()), err)
		}
		return handler.DBErrHandler(err)
	}
	// Update the fields of upModel as needed
	upModel.ObjID = Product.Id().Value()
	upModel.Name = Product.Name().Value()
	upModel.Price = int(Product.Price().Value())
	upModel.CategoryID = Product.Category().Id().Value()
	if _, updateErr := upModel.Update(ctx, tx, boil.Whitelist(
		models.ProductColumns.ObjID,
		models.ProductColumns.Name,
		models.ProductColumns.Price,
		models.ProductColumns.CategoryID,
	)); updateErr != nil {
		return handler.DBErrHandler(updateErr)
	}
	return nil
}

// DeleteById は指定されたIDの商品をデータベースから削除します。
//
// Parameters:
//   - ctx: コンテキスト
//   - tx: トランザクション
//   - id: 削除する商品のID
//
// Returns:
//   - error: 商品が存在しない場合やデータベースエラー
func (r *ProductRepositoryImpl) DeleteById(ctx context.Context, tx *sql.Tx, id *products.ProductId) error {
	condition := models.ProductWhere.ObjID.EQ(id.Value())
	delModel, err := models.Products(condition).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewCRUDError("NOT_FOUND", fmt.Sprintf("商品番号: %s は存在しないため、削除できませんでした。", id.Value()))
		}
		return handler.DBErrHandler(err)
	}
	if _, deleteErr := delModel.Delete(ctx, tx); deleteErr != nil {
		return handler.DBErrHandler(deleteErr)
	}
	return nil
}

// productHookFactory は指定されたフックタイプに応じた商品用フック関数を生成します。
//
// Parameters:
//   - hookType: フックのタイプ
//   - logger: ロガー
//
// Returns:
//   - func: SQLBoilerのフック関数
//
// Panics:
//   - hookTypeが想定外の値の場合
func productHookFactory(hookType boil.HookPoint, logger *slog.Logger) func(ctx context.Context, exec boil.ContextExecutor, product *models.Product) error {
	switch hookType {
	case boil.AfterInsertHook:
		return func(ctx context.Context, exec boil.ContextExecutor, product *models.Product) error {
			logger.InfoContext(ctx, "商品を登録しました。",
				slog.String("obj_id", product.ObjID),
				slog.String("name", product.Name),
				slog.Int("price", product.Price),
				slog.String("category_id", product.CategoryID),
			)
			return nil
		}
	case boil.AfterUpdateHook:
		return func(ctx context.Context, exec boil.ContextExecutor, product *models.Product) error {
			logger.InfoContext(ctx, "商品を変更しました。",
				slog.String("obj_id", product.ObjID),
				slog.String("name", product.Name),
				slog.Int("price", product.Price),
				slog.String("category_id", product.CategoryID),
			)
			return nil
		}
	case boil.AfterDeleteHook:
		return func(ctx context.Context, exec boil.ContextExecutor, product *models.Product) error {
			logger.InfoContext(ctx, "商品を削除しました。",
				slog.String("obj_id", product.ObjID),
				slog.String("name", product.Name),
				slog.Int("price", product.Price),
				slog.String("category_id", product.CategoryID),
			)
			return nil
		}
	default:
		panic("Invalid hookType")
	}
}

var _ products.ProductRepository = (*ProductRepositoryImpl)(nil)
