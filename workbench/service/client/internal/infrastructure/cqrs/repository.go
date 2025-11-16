package cqrs

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	command "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/command/v1"
	common "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/common/v1"
	query "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/domain/models"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/domain/repository"
)

// CQRSRepositoryImpl はCQRSRepositoryの実装
type CQRSRepositoryImpl struct {
	logger               *slog.Logger          // ロガー
	commandServiceClient *CommandServiceClient // Command Serviceクライアント
	queryServiceClient   *QueryServiceClient   // Query Serviceクライアント
}

// NewCQRSRepositoryImpl はCQRSRepositoryImplを生成します。
//
// Parameters:
//   - commandServiceClient: Command Serviceクライアント
//   - queryServiceClient: Query Serviceクライアント
//   - logger: ロガー
//
// Returns:
//   - repository.CQRSRepository: CQRSRepositoryの実装
func NewCQRSRepositoryImpl(commandServiceClient *CommandServiceClient, queryServiceClient *QueryServiceClient, logger *slog.Logger) *CQRSRepositoryImpl {
	return &CQRSRepositoryImpl{
		logger:               logger,
		commandServiceClient: commandServiceClient,
		queryServiceClient:   queryServiceClient,
	}
}

// newCategoryId はカテゴリIDのprotobuf値オブジェクトを生成します。
//
// Parameters:
//   - id: カテゴリID
//
// Returns:
//   - *common.CategoryId: CategoryId
func newCategoryId(id string) *common.CategoryId {
	v := &common.CategoryId{}
	v.SetValue(id)
	return v
}

// newCategoryName はカテゴリ名のprotobuf値オブジェクトを生成します。
//
// Parameters:
//   - name: カテゴリ名
//
// Returns:
//   - *common.CategoryName: CategoryName
func newCategoryName(name string) *common.CategoryName {
	v := &common.CategoryName{}
	v.SetValue(name)
	return v
}

// newProductId は商品IDのprotobuf値オブジェクトを生成します。
//
// Parameters:
//   - id: 商品ID
//
// Returns:
//   - *common.ProductId: ProductId
func newProductId(id string) *common.ProductId {
	v := &common.ProductId{}
	v.SetValue(id)
	return v
}

// newProductName は商品名のprotobuf値オブジェクトを生成します。
//
// Parameters:
//   - name: 商品名
//
// Returns:
//   - *common.ProductName: ProductName
func newProductName(name string) *common.ProductName {
	v := &common.ProductName{}
	v.SetValue(name)
	return v
}

// newProductPrice は価格のprotobuf値オブジェクトを生成します。
//
// Parameters:
//   - price: 価格
//
// Returns:
//   - *common.ProductPrice: ProductPrice
func newProductPrice(price uint32) *common.ProductPrice {
	v := &common.ProductPrice{}
	v.SetValue(int32(price))
	return v
}

// newCreateProductCategory は商品作成用カテゴリのprotobufオブジェクトを生成します。
//
// Parameters:
//   - category: カテゴリドメインモデル
//
// Returns:
//   - *command.CreateProductRequest_Product_Category: CreateProductRequest_Product_Category
func newCreateProductCategory(category *models.Category) *command.CreateProductRequest_Product_Category {
	c := &command.CreateProductRequest_Product_Category{}
	c.SetId(newCategoryId(category.Id()))
	c.SetName(newCategoryName(category.Name()))
	return c
}

// CreateCategory はカテゴリを作成します。
//
// Parameters:
//   - ctx: コンテキスト
//   - categoryName: カテゴリ名
//
// Returns:
//   - *models.Category: 作成されたカテゴリ
//   - error: エラー
func (r *CQRSRepositoryImpl) CreateCategory(ctx context.Context, categoryName string) (*models.Category, error) {
	req := &command.CreateCategoryRequest{}
	req.SetName(newCategoryName(categoryName))
	req.SetCrud(command.CRUD_CRUD_INSERT)

	resp, err := r.commandServiceClient.Category.CreateCategory(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return toModelCategory(resp.Msg.GetCategory()), nil
}

// UpdateCategory はカテゴリを更新します。
//
// Parameters:
//   - ctx: コンテキスト
//   - category: 更新するカテゴリ
//
// Returns:
//   - *models.Category: 更新されたカテゴリ
//   - error: エラー
func (r *CQRSRepositoryImpl) UpdateCategory(ctx context.Context, category *models.Category) (*models.Category, error) {
	c := &command.UpdateCategoryRequest_Category{}
	c.SetId(newCategoryId(category.Id()))
	c.SetName(newCategoryName(category.Name()))

	req := &command.UpdateCategoryRequest{}
	req.SetCategory(c)
	req.SetCrud(command.CRUD_CRUD_UPDATE)

	resp, err := r.commandServiceClient.Category.UpdateCategory(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return toModelCategory(resp.Msg.GetCategory()), nil
}

// DeleteCategory はカテゴリを削除します。
//
// Parameters:
//   - ctx: コンテキスト
//   - id: 削除するカテゴリID
//
// Returns:
//   - error: エラー
func (r *CQRSRepositoryImpl) DeleteCategory(ctx context.Context, id string) error {
	req := &command.DeleteCategoryRequest{}
	req.SetCategoryId(newCategoryId(id))
	req.SetCrud(command.CRUD_CRUD_DELETE)

	_, err := r.commandServiceClient.Category.DeleteCategory(ctx, connect.NewRequest(req))
	return err
}

// CategoryList はカテゴリ一覧を取得します。
//
// Parameters:
//   - ctx: コンテキスト
//
// Returns:
//   - []*models.Category: カテゴリ一覧
//   - error: エラー
func (r *CQRSRepositoryImpl) CategoryList(ctx context.Context) ([]*models.Category, error) {
	resp, err := r.queryServiceClient.Category.ListCategories(ctx, connect.NewRequest(&query.ListCategoriesRequest{}))
	if err != nil {
		return nil, err
	}
	return toModelCategories(resp.Msg.GetCategories()), nil
}

// CategoryById はIDでカテゴリを取得します。
//
// Parameters:
//   - ctx: コンテキスト
//   - id: カテゴリID
//
// Returns:
//   - *models.Category: カテゴリ
//   - error: エラー
func (r *CQRSRepositoryImpl) CategoryById(ctx context.Context, id string) (*models.Category, error) {
	req := &query.GetCategoryByIdRequest{}
	req.SetId(id)

	resp, err := r.queryServiceClient.Category.GetCategoryById(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return toModelCategory(resp.Msg.GetCategory()), nil
}

// CreateProduct は商品を作成します。
//
// Parameters:
//   - ctx: コンテキスト
//   - productName: 商品名
//   - productPrice: 価格
//   - category: カテゴリ
//
// Returns:
//   - *models.Product: 作成された商品
//   - error: エラー
func (r *CQRSRepositoryImpl) CreateProduct(ctx context.Context, productName string, productPrice uint32, category *models.Category) (*models.Product, error) {
	// Product nested message
	p := &command.CreateProductRequest_Product{}
	p.SetName(newProductName(productName))
	p.SetPrice(newProductPrice(productPrice))
	p.SetCategory(newCreateProductCategory(category))

	req := &command.CreateProductRequest{}
	req.SetProduct(p)
	req.SetCrud(command.CRUD_CRUD_INSERT)

	resp, err := r.commandServiceClient.Product.CreateProduct(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return toModelProduct(resp.Msg.GetProduct()), nil
}

// UpdateProduct は商品を更新します。
//
// Parameters:
//   - ctx: コンテキスト
//   - product: 更新する商品
//
// Returns:
//   - *models.Product: 更新された商品
//   - error: エラー
func (r *CQRSRepositoryImpl) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	p := &command.UpdateProductRequest_Product{}
	p.SetId(newProductId(product.Id()))
	p.SetName(newProductName(product.Name()))
	p.SetPrice(newProductPrice(product.Price()))
	p.SetCategoryId(newCategoryId(product.Category().Id()))

	req := &command.UpdateProductRequest{}
	req.SetProduct(p)
	req.SetCrud(command.CRUD_CRUD_UPDATE)

	resp, err := r.commandServiceClient.Product.UpdateProduct(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return toModelProduct(resp.Msg.GetProduct()), nil
}

// DeleteProduct は商品を削除します。
//
// Parameters:
//   - ctx: コンテキスト
//   - id: 削除する商品ID
//
// Returns:
//   - error: エラー
func (r *CQRSRepositoryImpl) DeleteProduct(ctx context.Context, id string) error {
	req := &command.DeleteProductRequest{}
	req.SetProductId(newProductId(id))

	_, err := r.commandServiceClient.Product.DeleteProduct(ctx, connect.NewRequest(req))
	return err
}

// ProductList は商品一覧を取得します。
//
// Parameters:
//   - ctx: コンテキスト
//
// Returns:
//   - []*models.Product: 商品一覧
//   - error: エラー
func (r *CQRSRepositoryImpl) ProductList(ctx context.Context) ([]*models.Product, error) {
	resp, err := r.queryServiceClient.Product.ListProducts(ctx, connect.NewRequest(&query.ListProductsRequest{}))
	if err != nil {
		return nil, err
	}
	return toModelProducts(resp.Msg.GetProducts()), nil
}

// StreamProducts はQuery ServiceのサーバーストリーミングRPCを呼び出し、
// 受信結果をchannelで返します。
//
// Parameters:
//   - ctx: コンテキスト
//
// Returns:
//   - <-chan *repository.StreamProductsResult: ストリーム結果
//   - error: エラー
func (r *CQRSRepositoryImpl) StreamProducts(ctx context.Context) (<-chan *repository.StreamProductsResult, error) {
	stream, err := r.queryServiceClient.Product.StreamProducts(ctx, connect.NewRequest(&query.StreamProductsRequest{}))
	if err != nil {
		return nil, err
	}

	ch := make(chan *repository.StreamProductsResult)
	go func() {
		defer func() {
			if closeErr := stream.Close(); closeErr != nil {
				r.logger.ErrorContext(ctx, "Failed to close stream", "error", closeErr)
				select {
				case ch <- &repository.StreamProductsResult{Err: closeErr}:
				case <-ctx.Done():
				default:
					r.logger.WarnContext(ctx, "Failed to send stream close error to channel as it is blocking", "error", closeErr)
				}
			}
			close(ch)
		}()
		for stream.Receive() {
			msg := stream.Msg()
			if msg == nil {
				continue
			}
			productProto := msg.GetProduct()
			if productProto == nil {
				continue
			}
			select {
			case ch <- &repository.StreamProductsResult{Product: toModelProduct(productProto)}:
			case <-ctx.Done():
				return
			}
		}
		if streamErr := stream.Err(); streamErr != nil {
			select {
			case ch <- &repository.StreamProductsResult{Err: streamErr}:
			case <-ctx.Done():
			}
		}
	}()
	return ch, nil
}

// ProductById はIDで商品を取得します。
//
// Parameters:
//   - ctx: コンテキスト
//   - id: 商品ID
//
// Returns:
//   - *models.Product: 商品
//   - error: エラー
func (r *CQRSRepositoryImpl) ProductById(ctx context.Context, id string) (*models.Product, error) {
	req := &query.GetProductByIdRequest{}
	req.SetId(id)

	resp, err := r.queryServiceClient.Product.GetProductById(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return toModelProduct(resp.Msg.GetProduct()), nil
}

// ProductByKeyword はキーワードで商品を検索します。
//
// Parameters:
//   - ctx: コンテキスト
//   - keyword: 検索キーワード
//
// Returns:
//   - []*models.Product: 検索結果の商品一覧
//   - error: エラー
func (r *CQRSRepositoryImpl) ProductByKeyword(ctx context.Context, keyword string) ([]*models.Product, error) {
	req := &query.SearchProductsByKeywordRequest{}
	req.SetKeyword(keyword)

	resp, err := r.queryServiceClient.Product.SearchProductsByKeyword(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return toModelProducts(resp.Msg.GetProducts()), nil
}

// toModelCategory はprotobufのCategoryをドメインモデルに変換します。
//
// Parameters:
//   - category: protobuf Category
//
// Returns:
//   - *models.Category: Categoryドメインモデル
func toModelCategory(category *common.Category) *models.Category {
	return models.NewCategory(category.GetId(), category.GetName())
}

// toModelCategories はprotobufのCategoryスライスをドメインモデルスライスに変換します。
//
// Parameters:
//   - categories: protobuf Categoryスライス
//
// Returns:
//   - []*models.Category: Categoryドメインモデルスライス
func toModelCategories(categories []*common.Category) []*models.Category {
	result := make([]*models.Category, len(categories))
	for i, category := range categories {
		result[i] = toModelCategory(category)
	}
	return result
}

// toModelProduct はprotobufのProductをドメインモデルに変換します。
//
// Parameters:
//   - product: protobuf Product
//
// Returns:
//   - *models.Product: Productドメインモデル
func toModelProduct(product *common.Product) *models.Product {
	return models.NewProduct(product.GetId(), product.GetName(), uint32(product.GetPrice()), toModelCategory(product.GetCategory()))
}

// toModelProducts はprotobufのProductスライスをドメインモデルスライスに変換します。
//
// Parameters:
//   - products: protobuf Productスライス
//
// Returns:
//   - []*models.Product: Productドメインモデルスライス
func toModelProducts(products []*common.Product) []*models.Product {
	result := make([]*models.Product, len(products))
	for i, product := range products {
		result[i] = toModelProduct(product)
	}
	return result
}

var _ repository.CQRSRepository = (*CQRSRepositoryImpl)(nil)
