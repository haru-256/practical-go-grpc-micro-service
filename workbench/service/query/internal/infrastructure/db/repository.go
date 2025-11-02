package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/domain/models"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/domain/repository"
	"gorm.io/gorm"
)

const (
	PRODUCT_ID_COLUMN   = "obj_id"
	PRODUCT_NAME_COLUMN = "name"
)

type ProductRepositoryImpl struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewProductRepositoryImpl(db *gorm.DB, logger *slog.Logger) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db, logger: logger}
}

func (r *ProductRepositoryImpl) List(ctx context.Context) ([]*models.Product, error) {
	products := []*Product{}
	if result := r.db.WithContext(ctx).Preload("Category").Find(&products); result.Error != nil {
		return nil, DBErrHandler(ctx, result.Error, r.logger)
	}

	return toProductModels(products), nil
}

func (r *ProductRepositoryImpl) FindById(ctx context.Context, id string) (*models.Product, error) {
	product := &Product{}
	if result := r.db.WithContext(ctx).Preload("Category").Where(fmt.Sprintf("%s = ?", PRODUCT_ID_COLUMN), id).First(product); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NewCRUDError("NOT_FOUND", fmt.Sprintf("商品ID: %s が見つかりませんでした", id))
		}
		return nil, DBErrHandler(ctx, result.Error, r.logger)
	}

	return toProductModel(product), nil
}

func (r *ProductRepositoryImpl) FindByNameLike(ctx context.Context, keyword string) ([]*models.Product, error) {
	if keyword == "" {
		return nil, errs.NewInternalError("INVALID_KEYWORD", "検索キーワードが空です")
	}

	products := []*Product{}
	likePattern := "%" + keyword + "%"
	if result := r.db.WithContext(ctx).Preload("Category").Where(fmt.Sprintf("%s LIKE ?", PRODUCT_NAME_COLUMN), likePattern).Find(&products); result.Error != nil {
		return nil, DBErrHandler(ctx, result.Error, r.logger)
	}

	return toProductModels(products), nil
}

func toProductModels(products []*Product) []*models.Product {
	results := make([]*models.Product, len(products))
	for i, p := range products {
		results[i] = toProductModel(p)
	}
	return results
}
func toProductModel(product *Product) *models.Product {
	category := models.NewCategory(product.Category.ObjId, product.Category.Name)
	return models.NewProduct(product.ObjId, product.Name, product.Price, category)
}

var _ repository.ProductRepository = (*ProductRepositoryImpl)(nil)

type CategoryRepositoryImpl struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewCategoryRepositoryImpl はCategoryRepositoryImplを生成します。
//
// Parameters:
//   - db: データベース接続
//
// Returns:
//   - *CategoryRepositoryImpl: CategoryRepositoryImplポインタ
func NewCategoryRepositoryImpl(db *gorm.DB, logger *slog.Logger) *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{db: db, logger: logger}
}

// List はすべてのカテゴリを取得します。
//
// Parameters:
//   - ctx: コンテキスト
//
// Returns:
//   - []*models.Category: カテゴリリスト
//   - error: エラー
func (r *CategoryRepositoryImpl) List(ctx context.Context) ([]*models.Category, error) {
	categories := []*Category{}
	if result := r.db.WithContext(ctx).Find(&categories); result.Error != nil {
		return nil, DBErrHandler(ctx, result.Error, r.logger)
	}

	return toCategoryModels(categories), nil
}

// FindById はカテゴリIDでカテゴリを検索します。
//
// Parameters:
//   - ctx: コンテキスト
//   - id: カテゴリID
//
// Returns:
//   - *models.Category: カテゴリ
//   - error: エラー
func (r *CategoryRepositoryImpl) FindById(ctx context.Context, id string) (*models.Category, error) {
	category := &Category{}
	if result := r.db.WithContext(ctx).Where("obj_id = ?", id).First(category); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NewCRUDError("NOT_FOUND", fmt.Sprintf("カテゴリID: %s が見つかりませんでした", id))
		}
		return nil, DBErrHandler(ctx, result.Error, r.logger)
	}

	return toCategoryModel(category), nil
}

func toCategoryModels(categories []*Category) []*models.Category {
	results := make([]*models.Category, len(categories))
	for i, c := range categories {
		results[i] = toCategoryModel(c)
	}
	return results
}

func toCategoryModel(category *Category) *models.Category {
	return models.NewCategory(category.ObjId, category.Name)
}

var _ repository.CategoryRepository = (*CategoryRepositoryImpl)(nil)
