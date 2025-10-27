package dto

import (
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
)

// CategoryDTO はカテゴリデータのDTOです。
type CategoryDTO struct {
	Id   string // カテゴリID
	Name string // カテゴリ名
}

// CreateCategoryDTO はカテゴリの新規作成時に使用するDTOです。
type CreateCategoryDTO struct {
	Name string // カテゴリ名
}

// UpdateCategoryDTO はカテゴリの更新時に使用するDTOです。
type UpdateCategoryDTO struct {
	Id   string // カテゴリID
	Name string // カテゴリ名
}

// DeleteCategoryDTO はカテゴリの削除時に使用するDTOです。
type DeleteCategoryDTO struct {
	Id string // カテゴリID
}

// ProductDTO は商品データのDTOです。
type ProductDTO struct {
	Id       string       // 商品ID
	Name     string       // 商品名
	Category *CategoryDTO // 商品カテゴリ
	Price    uint32       // 単価
}

// CreateProductDTO は商品の新規作成時に使用するDTOです。
type CreateProductDTO struct {
	Name     string       // 商品名
	Price    uint32       // 単価
	Category *CategoryDTO // 既存カテゴリ情報
}

// UpdateProductDTO は商品の更新時に使用するDTOです。
type UpdateProductDTO struct {
	Id         string // 商品ID
	Name       string // 商品名
	Price      uint32 // 単価
	CategoryId string // 商品カテゴリID
}

// DeleteProductDTO は商品の削除時に使用するDTOです。
type DeleteProductDTO struct {
	Id string // 商品ID
}

// NewCategoryDTOFromEntity はドメインエンティティからDTOを生成します。
//
// Parameters:
//   - category: 変換元のカテゴリエンティティ
//
// Returns:
//   - *CategoryDTO: プレゼンテーション層で使用するDTO
func NewCategoryDTOFromEntity(category *categories.Category) *CategoryDTO {
	return &CategoryDTO{
		Id:   category.Id().Value(),
		Name: category.Name().Value(),
	}
}

// CategoryFromDTO はDTOからドメインエンティティを再構築します。
//
// Parameters:
//   - dto: 変換元のDTO
//
// Returns:
//   - *categories.Category: 再構築されたドメインエンティティ
//   - error: 変換エラー
func CategoryFromDTO(dto *CategoryDTO) (*categories.Category, error) {
	id, err := categories.NewCategoryId(dto.Id)
	if err != nil {
		return nil, err
	}
	name, err := categories.NewCategoryName(dto.Name)
	if err != nil {
		return nil, err
	}
	return categories.BuildCategory(id, name)
}

// CategoryFromCreateDTO は新規作成用DTOからドメインエンティティを生成します。
//
// Parameters:
//   - dto: 変換元のDTO
//
// Returns:
//   - *categories.Category: 生成されたドメインエンティティ
//   - error: 変換エラー
func CategoryFromCreateDTO(dto *CreateCategoryDTO) (*categories.Category, error) {
	name, err := categories.NewCategoryName(dto.Name)
	if err != nil {
		return nil, err
	}
	return categories.NewCategory(name)
}

// CategoryFromUpdateDTO は更新用DTOからドメインエンティティを再構築します。
//
// Parameters:
//   - dto: 変換元のDTO
//
// Returns:
//   - *categories.Category: 再構築されたドメインエンティティ
//   - error: 変換エラー
func CategoryFromUpdateDTO(dto *UpdateCategoryDTO) (*categories.Category, error) {
	id, err := categories.NewCategoryId(dto.Id)
	if err != nil {
		return nil, err
	}
	name, err := categories.NewCategoryName(dto.Name)
	if err != nil {
		return nil, err
	}
	return categories.BuildCategory(id, name)
}

// CategoryIdFromDeleteDTO は削除用DTOからカテゴリIDを取得します。
//
// Parameters:
//   - dto: 変換元のDTO
//
// Returns:
//   - *categories.CategoryId: カテゴリID
//   - error: 変換エラー
func CategoryIdFromDeleteDTO(dto *DeleteCategoryDTO) (*categories.CategoryId, error) {
	return categories.NewCategoryId(dto.Id)
}

// NewProductDTOFromEntity はドメインエンティティからDTOを生成します。
//
// Parameters:
//   - product: 変換元の商品エンティティ
//
// Returns:
//   - *ProductDTO: プレゼンテーション層で使用するDTO
func NewProductDTOFromEntity(product *products.Product) *ProductDTO {
	return &ProductDTO{
		Id:   product.Id().Value(),
		Name: product.Name().Value(),
		Category: &CategoryDTO{
			Id:   product.Category().Id().Value(),
			Name: product.Category().Name().Value(),
		},
		Price: product.Price().Value(),
	}
}

// ProductFromDTO はDTOからドメインエンティティを再構築します。
//
// Parameters:
//   - dto: 変換元のDTO
//
// Returns:
//   - *products.Product: 再構築されたドメインエンティティ
//   - error: 変換エラー
func ProductFromDTO(dto *ProductDTO) (*products.Product, error) {
	id, err := products.NewProductId(dto.Id)
	if err != nil {
		return nil, err
	}
	name, err := products.NewProductName(dto.Name)
	if err != nil {
		return nil, err
	}
	price, err := products.NewProductPrice(dto.Price)
	if err != nil {
		return nil, err
	}
	category, err := CategoryFromDTO(dto.Category)
	if err != nil {
		return nil, err
	}
	return products.BuildProduct(id, name, price, category)
}

// ProductFromCreateDTO は新規作成用DTOからドメインエンティティを生成します。
//
// Parameters:
//   - dto: 変換元のDTO
//
// Returns:
//   - *products.Product: 生成されたドメインエンティティ
//   - error: 変換エラー
func ProductFromCreateDTO(dto *CreateProductDTO) (*products.Product, error) {
	name, err := products.NewProductName(dto.Name)
	if err != nil {
		return nil, err
	}
	price, err := products.NewProductPrice(dto.Price)
	if err != nil {
		return nil, err
	}
	category, err := CategoryFromDTO(dto.Category)
	if err != nil {
		return nil, err
	}
	return products.NewProduct(name, price, category)
}

// ProductFromUpdateDTO は更新用DTOからドメインエンティティを再構築します。
//
// Parameters:
//   - dto: 変換元のDTO
//   - categoryName: カテゴリ名
//
// Returns:
//   - *products.Product: 再構築されたドメインエンティティ
//   - error: 変換エラー
func ProductFromUpdateDTO(dto *UpdateProductDTO, categoryName *categories.CategoryName) (*products.Product, error) {
	id, err := products.NewProductId(dto.Id)
	if err != nil {
		return nil, err
	}
	name, err := products.NewProductName(dto.Name)
	if err != nil {
		return nil, err
	}
	price, err := products.NewProductPrice(dto.Price)
	if err != nil {
		return nil, err
	}
	categoryDTO := &CategoryDTO{
		Id:   dto.CategoryId,
		Name: categoryName.Value(),
	}
	category, err := CategoryFromDTO(categoryDTO)
	if err != nil {
		return nil, err
	}
	return products.BuildProduct(id, name, price, category)
}

// ProductIdFromDeleteDTO は削除用DTOから商品IDを取得します。
//
// Parameters:
//   - dto: 変換元のDTO
//
// Returns:
//   - *products.ProductId: 商品ID
//   - error: 変換エラー
func ProductIdFromDeleteDTO(dto *DeleteProductDTO) (*products.ProductId, error) {
	return products.NewProductId(dto.Id)
}
