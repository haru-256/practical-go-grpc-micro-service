package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	common "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/common/v1"
	query "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1"
	queryconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1/queryv1connect"
	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/domain/models"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/domain/repository"
)

// CategoryServiceHandlerImpl はCategoryServiceのgRPCハンドラ実装です。
type CategoryServiceHandlerImpl struct {
	logger *slog.Logger                  // ロガー
	repo   repository.CategoryRepository // カテゴリリポジトリ
	queryconnect.UnimplementedCategoryServiceHandler
}

// NewCategoryServiceHandlerImpl はCategoryServiceHandlerImplを生成します。
//
// Parameters:
//   - logger: ロガー
//   - validator: バリデータ
//   - repo: カテゴリリポジトリ
//
// Returns:
//   - *CategoryServiceHandlerImpl: ハンドラインスタンス
//   - error: エラー
func NewCategoryServiceHandlerImpl(logger *slog.Logger, repo repository.CategoryRepository) (*CategoryServiceHandlerImpl, error) {
	return &CategoryServiceHandlerImpl{
		logger: logger,
		repo:   repo,
	}, nil
}

// ListCategories はカテゴリ一覧を取得します。
//
// Parameters:
//   - ctx: コンテキスト
//   - req: リクエスト
//
// Returns:
//   - *connect.Response[query.ListCategoriesResponse]: レスポンス
//   - error: エラー
func (h *CategoryServiceHandlerImpl) ListCategories(ctx context.Context, req *connect.Request[query.ListCategoriesRequest]) (*connect.Response[query.ListCategoriesResponse], error) {
	// カテゴリ一覧を取得
	categories, err := h.repo.List(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to list categories", "error", err)
		return nil, handleError(err, "failed to list categories")
	}

	// レスポンス生成
	res := &query.ListCategoriesResponse{}
	res.SetCategories(toCategoriesProto(categories))

	return connect.NewResponse(res), nil
}

// GetCategoryById はカテゴリIDでカテゴリを取得します。
//
// Parameters:
//   - ctx: コンテキスト
//   - req: リクエスト
//
// Returns:
//   - *connect.Response[query.GetCategoryByIdResponse]: レスポンス
//   - error: エラー
func (h *CategoryServiceHandlerImpl) GetCategoryById(ctx context.Context, req *connect.Request[query.GetCategoryByIdRequest]) (*connect.Response[query.GetCategoryByIdResponse], error) {
	// カテゴリを取得
	category, err := h.repo.FindById(ctx, req.Msg.GetId())
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to get category by id", "error", err, "id", req.Msg.GetId())
		return nil, handleError(err, "failed to get category by id")
	}

	// レスポンス生成
	res := &query.GetCategoryByIdResponse{}
	res.SetCategory(toCategoryProto(category))

	return connect.NewResponse(res), nil
}

// ProductServiceHandlerImpl はProductServiceのgRPCハンドラ実装です。
type ProductServiceHandlerImpl struct {
	logger *slog.Logger                 // ロガー
	repo   repository.ProductRepository // 商品リポジトリ
	queryconnect.UnimplementedProductServiceHandler
}

// NewProductServiceHandlerImpl はProductServiceHandlerImplを生成します。
//
// Parameters:
//   - logger: ロガー
//   - validator: バリデータ
//   - repo: 商品リポジトリ
//
// Returns:
//   - *ProductServiceHandlerImpl: ハンドラインスタンス
//   - error: エラー
func NewProductServiceHandlerImpl(logger *slog.Logger, repo repository.ProductRepository) (*ProductServiceHandlerImpl, error) {
	return &ProductServiceHandlerImpl{
		logger: logger,
		repo:   repo,
	}, nil
}

// ListProducts は商品一覧を取得します。
//
// Parameters:
//   - ctx: コンテキスト
//   - req: リクエスト
//
// Returns:
//   - *connect.Response[query.ListProductsResponse]: レスポンス
//   - error: エラー
func (h *ProductServiceHandlerImpl) ListProducts(ctx context.Context, req *connect.Request[query.ListProductsRequest]) (*connect.Response[query.ListProductsResponse], error) {
	// 商品取得
	products, err := h.repo.List(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to list products", "error", err)
		return nil, handleError(err, "failed to list products")
	}

	// レスポンス生成
	res := &query.ListProductsResponse{}
	res.SetProducts(toProductsProto(products))

	return connect.NewResponse(res), nil
}

// StreamProducts は商品一覧をサーバーストリームで返します。
//
// Parameters:
//   - ctx: コンテキスト
//   - req: リクエスト
//   - stream: ストリームレスポンス
//
// Returns:
//   - error: エラー
func (h *ProductServiceHandlerImpl) StreamProducts(ctx context.Context, req *connect.Request[query.StreamProductsRequest], stream *connect.ServerStream[query.StreamProductsResponse]) error {
	if products, err := h.repo.List(ctx); err != nil {
		h.logger.ErrorContext(ctx, "Failed to list products", "error", err)
		return handleError(err, "failed to list products")
	} else {
		for _, product := range products {
			res := &query.StreamProductsResponse{}
			res.SetProduct(toProductProto(product))
			if sendErr := stream.Send(res); sendErr != nil {
				h.logger.ErrorContext(ctx, "Failed to send product in stream", "error", sendErr, "product_id", product.Id())
				return handleError(sendErr, "failed to send product in stream")
			}
		}
	}
	return nil
}

// GetProductById は商品IDで商品を取得します。
//
// Parameters:
//   - ctx: コンテキスト
//   - req: リクエスト
//
// Returns:
//   - *connect.Response[query.GetProductByIdResponse]: レスポンス
//   - error: エラー
func (h *ProductServiceHandlerImpl) GetProductById(ctx context.Context, req *connect.Request[query.GetProductByIdRequest]) (*connect.Response[query.GetProductByIdResponse], error) {
	// 商品を取得
	product, err := h.repo.FindById(ctx, req.Msg.GetId())
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to get product by id", "error", err, "id", req.Msg.GetId())
		return nil, handleError(err, "failed to get product by id")
	}

	// レスポンス生成
	res := &query.GetProductByIdResponse{}
	res.SetProduct(toProductProto(product))

	return connect.NewResponse(res), nil
}

// SearchProductsByKeyword は商品名のキーワードで商品を検索します。
//
// Parameters:
//   - ctx: コンテキスト
//   - req: リクエスト
//
// Returns:
//   - *connect.Response[query.SearchProductsByKeywordResponse]: レスポンス
//   - error: エラー
func (h *ProductServiceHandlerImpl) SearchProductsByKeyword(ctx context.Context, req *connect.Request[query.SearchProductsByKeywordRequest]) (*connect.Response[query.SearchProductsByKeywordResponse], error) {
	// 商品取得
	products, err := h.repo.FindByNameLike(ctx, req.Msg.GetKeyword())
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to search products by keyword", "error", err, "keyword", req.Msg.GetKeyword())
		return nil, handleError(err, "failed to search products by keyword")
	}

	// レスポンス生成
	res := &query.SearchProductsByKeywordResponse{}
	res.SetProducts(toProductsProto(products))

	return connect.NewResponse(res), nil
}

// toCategoryProto はドメインモデルのCategoryをprotobufのCategoryに変換します。
//
// Parameters:
//   - category: ドメインモデルのCategory
//
// Returns:
//   - *common.Category: protobufのCategory
func toCategoryProto(category *models.Category) *common.Category {
	c := &common.Category{}
	c.SetId(category.Id())
	c.SetName(category.Name())
	return c
}

// toCategoriesProto はドメインモデルのCategoryスライスをprotobufのCategoryスライスに変換します。
//
// Parameters:
//   - categories: ドメインモデルのCategoryスライス
//
// Returns:
//   - []*common.Category: protobufのCategoryスライス
func toCategoriesProto(categories []*models.Category) []*common.Category {
	result := make([]*common.Category, len(categories))
	for i, category := range categories {
		result[i] = toCategoryProto(category)
	}
	return result
}

// toProductProto はドメインモデルのProductをprotobufのProductに変換します。
//
// Parameters:
//   - product: ドメインモデルのProduct
//
// Returns:
//   - *common.Product: protobufのProduct
func toProductProto(product *models.Product) *common.Product {
	p := &common.Product{}
	p.SetId(product.Id())
	p.SetName(product.Name())
	p.SetPrice(int32(product.Price()))
	if product.Category() != nil {
		p.SetCategory(toCategoryProto(product.Category()))
	}
	return p
}

// toProductsProto はドメインモデルのProductスライスをprotobufのProductスライスに変換します。
//
// Parameters:
//   - products: ドメインモデルのProductスライス
//
// Returns:
//   - []*common.Product: protobufのProductスライス
func toProductsProto(products []*models.Product) []*common.Product {
	result := make([]*common.Product, len(products))
	for i, product := range products {
		result[i] = toProductProto(product)
	}
	return result
}

// handleError はドメインエラーを適切なgRPCエラーに変換します。
//
// Parameters:
//   - err: ドメインエラー
//   - operation: 操作名
//
// Returns:
//   - error: gRPCエラー
func handleError(err error, operation string) error {
	var crudErr *errs.CRUDError
	var internalErr *errs.InternalError

	if errors.As(err, &crudErr) {
		switch crudErr.Code {
		case "NOT_FOUND":
			return connect.NewError(connect.CodeNotFound, fmt.Errorf("%s: %w", operation, err))
		case "ALREADY_EXISTS":
			return connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("%s: %w", operation, err))
		default:
			return connect.NewError(connect.CodeInternal, fmt.Errorf("%s: %w", operation, err))
		}
	}

	if errors.As(err, &internalErr) {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("%s: %w", operation, err))
	}

	// 不明なエラー
	return connect.NewError(connect.CodeInternal, fmt.Errorf("%s: %w", operation, err))
}
