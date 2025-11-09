package server

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/domain/models"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/domain/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/presentation/dto"
	"github.com/labstack/echo/v4"
)

// CustomValidator はEchoのバリデータインターフェースを実装する構造体
type CustomValidator struct {
	validator *validator.Validate // validator/v10のバリデータ
}

// NewCustomValidator はCustomValidatorを生成します。
//
// Returns:
//   - *CustomValidator: CustomValidatorのインスタンス
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate はリクエストボディのバリデーションを行います。
//
// Parameters:
//   - i: バリデーション対象の構造体
//
// Returns:
//   - error: バリデーションエラー（エラーがない場合はnil）
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// CQRSServiceHandler はCQRSサービスのHTTPハンドラ
type CQRSServiceHandler struct {
	logger *slog.Logger              // ロガー
	repo   repository.CQRSRepository // CQRSリポジトリ
}

// NewCQRSServiceHandler はCQRSServiceHandlerを生成します。
//
// Parameters:
//   - logger: ロガー
//   - repo: CQRSリポジトリ
//
// Returns:
//   - *CQRSServiceHandler: CQRSServiceHandlerのインスタンス
func NewCQRSServiceHandler(logger *slog.Logger, repo repository.CQRSRepository) *CQRSServiceHandler {
	return &CQRSServiceHandler{
		logger: logger,
		repo:   repo,
	}
}

// CreateCategory はカテゴリを登録します。
// @tags Category
// @Summary カテゴリ登録
// @Description カテゴリを登録します。
// @ID create-category
// @Accept application/json
// @Produce application/json
// @Param request body dto.CreateCategoryRequest true "カテゴリ情報"
// @Success 201 {object} dto.CreateCategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories [post]
func (h *CQRSServiceHandler) CreateCategory(c echo.Context) error {
	req := new(dto.CreateCategoryRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// バリデーション
	// TODO: Interceptorで共通化する
	if err := c.Validate(req); err != nil {
		h.logger.Warn("Validation failed", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	category, err := h.repo.CreateCategory(c.Request().Context(), req.Name)
	if err != nil {
		h.logger.Error("Failed to create category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create category").SetInternal(err)
	}

	resp := dto.CreateCategoryResponse{
		Category: &dto.Category{
			Id:   category.Id(),
			Name: category.Name(),
		},
	}
	return c.JSON(http.StatusCreated, resp)
}

// CategoryList はカテゴリの一覧を取得します。
// @tags Category
// @Summary カテゴリ一覧取得
// @Description カテゴリの一覧を取得します。
// @ID list-categories
// @Produce application/json
// @Success 200 {object} dto.CategoryListResponse
// @Failure 500 {object} map[string]string
// @Router /categories [get]
func (h *CQRSServiceHandler) CategoryList(c echo.Context) error {
	categories, err := h.repo.CategoryList(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to list categories", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list categories").SetInternal(err)
	}

	var resp dto.CategoryListResponse
	for _, category := range categories {
		resp.Categories = append(resp.Categories, &dto.Category{
			Id:   category.Id(),
			Name: category.Name(),
		})
	}
	return c.JSON(http.StatusOK, resp)
}

// UpdateCategory はカテゴリを更新します。
// @tags Category
// @Summary カテゴリ更新
// @Description カテゴリを更新します。
// @ID update-category
// @Accept application/json
// @Produce application/json
// @Param id path string true "カテゴリID"
// @Param request body dto.UpdateCategoryRequest true "カテゴリ情報"
// @Success 200 {object} dto.UpdateCategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [put]
func (h *CQRSServiceHandler) UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	req := new(dto.UpdateCategoryRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.logger.Warn("Validation failed", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	category := models.NewCategory(id, req.Name)

	updated, err := h.repo.UpdateCategory(c.Request().Context(), category)
	if err != nil {
		h.logger.Error("Failed to update category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update category").SetInternal(err)
	}

	resp := dto.UpdateCategoryResponse{
		Category: &dto.Category{
			Id:   updated.Id(),
			Name: updated.Name(),
		},
	}
	return c.JSON(http.StatusOK, resp)
}

// DeleteCategory はカテゴリを削除します。
// @tags Category
// @Summary カテゴリ削除
// @Description カテゴリを削除します。
// @ID delete-category
// @Param id path string true "カテゴリID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [delete]
func (h *CQRSServiceHandler) DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	if err := h.repo.DeleteCategory(c.Request().Context(), id); err != nil {
		h.logger.Error("Failed to delete category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete category").SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// CategoryById はIDでカテゴリを取得します。
// @tags Category
// @Summary カテゴリ取得
// @Description IDでカテゴリを取得します。
// @ID get-category-by-id
// @Produce application/json
// @Param id path string true "カテゴリID"
// @Success 200 {object} dto.CategoryByIdResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [get]
func (h *CQRSServiceHandler) CategoryById(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	category, err := h.repo.CategoryById(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get category").SetInternal(err)
	}

	resp := dto.CategoryByIdResponse{
		Category: &dto.Category{
			Id:   category.Id(),
			Name: category.Name(),
		},
	}
	return c.JSON(http.StatusOK, resp)
}

// CreateProduct は商品を作成します。
// @tags Product
// @Summary 商品登録
// @Description 商品を登録します。
// @ID create-product
// @Accept application/json
// @Produce application/json
// @Param request body dto.CreateProductRequest true "商品情報"
// @Success 201 {object} dto.CreateProductResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func (h *CQRSServiceHandler) CreateProduct(c echo.Context) error {
	req := new(dto.CreateProductRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.logger.Warn("Validation failed", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	category := models.NewCategory(req.Category.Id, req.Category.Name)

	// FIXME: category nameは不要なはずなのに要求している
	product, err := h.repo.CreateProduct(c.Request().Context(), req.Name, req.Price, category)
	if err != nil {
		h.logger.Error("Failed to create product", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create product").SetInternal(err)
	}

	resp := dto.CreateProductResponse{
		Product: &dto.Product{
			Id:    product.Id(),
			Name:  product.Name(),
			Price: product.Price(),
			Category: &dto.Category{
				Id:   product.Category().Id(),
				Name: product.Category().Name(),
			},
		},
	}
	return c.JSON(http.StatusCreated, resp)
}

// UpdateProduct は商品を更新します。
// @tags Product
// @Summary 商品更新
// @Description 商品を更新します。
// @ID update-product
// @Accept application/json
// @Produce application/json
// @Param id path string true "商品ID"
// @Param request body dto.UpdateProductRequest true "商品情報"
// @Success 200 {object} dto.UpdateProductResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [put]
func (h *CQRSServiceHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	req := new(dto.UpdateProductRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.logger.Warn("Validation failed", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// FIXME: category nameは不要なはずなのに要求している
	category := models.NewCategory(req.Category.Id, req.Category.Name)
	product := models.NewProduct(id, req.Name, req.Price, category)

	updated, err := h.repo.UpdateProduct(c.Request().Context(), product)
	if err != nil {
		h.logger.Error("Failed to update product", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update product").SetInternal(err)
	}

	resp := dto.UpdateProductResponse{
		Product: &dto.Product{
			Id:    updated.Id(),
			Name:  updated.Name(),
			Price: updated.Price(),
			Category: &dto.Category{
				Id:   updated.Category().Id(),
				Name: updated.Category().Name(),
			},
		},
	}
	return c.JSON(http.StatusOK, resp)
}

// DeleteProduct は商品を削除します。
// @tags Product
// @Summary 商品削除
// @Description 商品を削除します。
// @ID delete-product
// @Param id path string true "商品ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [delete]
func (h *CQRSServiceHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	if err := h.repo.DeleteProduct(c.Request().Context(), id); err != nil {
		h.logger.Error("Failed to delete product", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete product").SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ProductList は商品一覧を取得します。keywordパラメータがある場合は検索を行います。
// @tags Product
// @Summary 商品一覧取得・検索
// @Description 商品一覧を取得します。keywordパラメータを指定すると検索を行います。
// @ID list-products
// @Produce application/json
// @Param keyword query string false "検索キーワード"
// @Success 200 {object} dto.ProductListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *CQRSServiceHandler) ProductList(c echo.Context) error {
	keyword := c.QueryParam("keyword")

	// keywordパラメータがある場合は検索
	if keyword != "" {
		return h.ProductByKeyword(c)
	}

	// keywordパラメータがない場合は一覧取得
	products, err := h.repo.ProductList(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to list products", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list products").SetInternal(err)
	}

	var resp dto.ProductListResponse
	for _, product := range products {
		resp.Products = append(resp.Products, &dto.Product{
			Id:    product.Id(),
			Name:  product.Name(),
			Price: product.Price(),
			Category: &dto.Category{
				Id:   product.Category().Id(),
				Name: product.Category().Name(),
			},
		})
	}
	return c.JSON(http.StatusOK, resp)
}

// ProductById はIDで商品を取得します。
// @tags Product
// @Summary 商品取得
// @Description IDで商品を取得します。
// @ID get-product-by-id
// @Produce application/json
// @Param id path string true "商品ID"
// @Success 200 {object} dto.ProductByIdResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [get]
func (h *CQRSServiceHandler) ProductById(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	product, err := h.repo.ProductById(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get product", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get product").SetInternal(err)
	}

	resp := dto.ProductByIdResponse{
		Product: &dto.Product{
			Id:    product.Id(),
			Name:  product.Name(),
			Price: product.Price(),
			Category: &dto.Category{
				Id:   product.Category().Id(),
				Name: product.Category().Name(),
			},
		},
	}
	return c.JSON(http.StatusOK, resp)
}

// ProductByKeyword はキーワードで商品を検索します。
// このメソッドはProductListから内部的に呼び出されます。
func (h *CQRSServiceHandler) ProductByKeyword(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	if keyword == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "keyword is required")
	}

	products, err := h.repo.ProductByKeyword(c.Request().Context(), keyword)
	if err != nil {
		h.logger.Error("Failed to search products", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to search products").SetInternal(err)
	}

	var resp dto.ProductByKeywordResponse
	for _, product := range products {
		resp.Products = append(resp.Products, &dto.Product{
			Id:    product.Id(),
			Name:  product.Name(),
			Price: product.Price(),
			Category: &dto.Category{
				Id:   product.Category().Id(),
				Name: product.Category().Name(),
			},
		})
	}
	return c.JSON(http.StatusOK, resp)
}
