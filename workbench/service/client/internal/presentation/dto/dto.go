package dto

// Category はカテゴリ情報を表すDTO
type Category struct {
	Id   string `json:"id" validate:"required,uuid4"`          // カテゴリID
	Name string `json:"name" validate:"required,min=1,max=20"` // カテゴリ名
}

// CreateCategoryRequest はカテゴリ作成リクエスト
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=20"` // カテゴリ名
}

// CreateCategoryResponse はカテゴリ作成レスポンス
type CreateCategoryResponse struct {
	Category *Category `json:"category"` // 作成されたカテゴリ情報
}

// CategoryListResponse はカテゴリ一覧レスポンス
type CategoryListResponse struct {
	Categories []*Category `json:"categories"` // カテゴリ一覧
}

// UpdateCategoryRequest はカテゴリ更新リクエスト
type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=20"` // カテゴリ名
}

// UpdateCategoryResponse はカテゴリ更新レスポンス
type UpdateCategoryResponse struct {
	Category *Category `json:"category"` // 更新されたカテゴリ情報
}

// CategoryByIdResponse はカテゴリ取得レスポンス
type CategoryByIdResponse struct {
	Category *Category `json:"category"` // カテゴリ情報
}

// Product は商品情報を表すDTO
type Product struct {
	Id       string    `json:"id"`       // 商品ID
	Name     string    `json:"name"`     // 商品名
	Price    uint32    `json:"price"`    // 価格
	Category *Category `json:"category"` // カテゴリ情報
}

// CreateProductRequest は商品作成リクエスト
type CreateProductRequest struct {
	Name     string    `json:"name" validate:"required,min=1,max=100"` // 商品名
	Price    uint32    `json:"price" validate:"required,min=1"`        // 価格
	Category *Category `json:"category" validate:"required"`           // カテゴリ情報
}

// CreateProductResponse は商品作成レスポンス
type CreateProductResponse struct {
	Product *Product `json:"product"` // 作成された商品情報
}

// UpdateProductRequest は商品更新リクエスト
type UpdateProductRequest struct {
	Name     string    `json:"name" validate:"required,min=1,max=100"` // 商品名
	Price    uint32    `json:"price" validate:"required,min=1"`        // 価格
	Category *Category `json:"category" validate:"required"`           // カテゴリ情報
}

// UpdateProductResponse は商品更新レスポンス
type UpdateProductResponse struct {
	Product *Product `json:"product"` // 更新された商品情報
}

// ProductListResponse は商品一覧レスポンス
type ProductListResponse struct {
	Products []*Product `json:"products"` // 商品一覧
}

// ProductByIdResponse は商品取得レスポンス
type ProductByIdResponse struct {
	Product *Product `json:"product"` // 商品情報
}

// ProductByKeywordResponse は商品検索レスポンス
type ProductByKeywordResponse struct {
	Products []*Product `json:"products"` // 検索結果の商品一覧
}
