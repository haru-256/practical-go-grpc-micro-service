package server_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/domain/models"
	mock_repository "github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/mock/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/presentation/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/presentation/server"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCQRSServiceHandler_CreateCategory(t *testing.T) {
	t.Run("正常系: カテゴリを作成できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":"TestCategory"}`
		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// モックの設定
		expectedCategory := models.NewCategory("cat-123", "TestCategory")
		mockRepo.EXPECT().
			CreateCategory(gomock.Any(), "TestCategory").
			Return(expectedCategory, nil)

		// Act
		err := handler.CreateCategory(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response dto.CreateCategoryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "cat-123", response.Category.Id)
		assert.Equal(t, "TestCategory", response.Category.Name)
	})

	t.Run("異常系: リクエストボディが不正", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `invalid json`
		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.CreateCategory(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("異常系: バリデーションエラー（nameが空）", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":""}`
		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.CreateCategory(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("異常系: バリデーションエラー（nameが長すぎる）", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":"ThisIsAVeryLongCategoryNameThatExceedsTwentyCharacters"}`
		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.CreateCategory(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("異常系: リポジトリエラー", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":"TestCategory"}`
		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// モックの設定: エラーを返す
		mockRepo.EXPECT().
			CreateCategory(gomock.Any(), "TestCategory").
			Return(nil, errors.New("database error"))

		// Act
		err := handler.CreateCategory(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)
	})

	t.Run("正常系: 境界値テスト（1文字）", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":"A"}`
		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// モックの設定
		expectedCategory := models.NewCategory("cat-123", "A")
		mockRepo.EXPECT().
			CreateCategory(gomock.Any(), "A").
			Return(expectedCategory, nil)

		// Act
		err := handler.CreateCategory(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("正常系: 境界値テスト（20文字）", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		categoryName := "12345678901234567890" // 20文字
		requestBody := `{"name":"` + categoryName + `"}`
		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// モックの設定
		expectedCategory := models.NewCategory("cat-123", categoryName)
		mockRepo.EXPECT().
			CreateCategory(gomock.Any(), categoryName).
			Return(expectedCategory, nil)

		// Act
		err := handler.CreateCategory(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestCQRSServiceHandler_CategoryList(t *testing.T) {
	t.Run("正常系: カテゴリ一覧を取得できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// モックの設定
		expectedCategories := []*models.Category{
			models.NewCategory("cat-1", "Category1"),
			models.NewCategory("cat-2", "Category2"),
		}
		mockRepo.EXPECT().
			CategoryList(gomock.Any()).
			Return(expectedCategories, nil)

		// Act
		err := handler.CategoryList(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.CategoryListResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response.Categories, 2)
		assert.Equal(t, "cat-1", response.Categories[0].Id)
		assert.Equal(t, "Category1", response.Categories[0].Name)
	})

	t.Run("異常系: リポジトリエラー", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.EXPECT().
			CategoryList(gomock.Any()).
			Return(nil, errors.New("database error"))

		// Act
		err := handler.CategoryList(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)
	})
}

func TestCQRSServiceHandler_UpdateCategory(t *testing.T) {
	t.Run("正常系: カテゴリを更新できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":"UpdatedCategory"}`
		req := httptest.NewRequest(http.MethodPut, "/categories/cat-123", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/categories/:id")
		c.SetParamNames("id")
		c.SetParamValues("cat-123")

		// モックの設定
		expectedCategory := models.NewCategory("cat-123", "UpdatedCategory")
		mockRepo.EXPECT().
			UpdateCategory(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx interface{}, category *models.Category) (*models.Category, error) {
				assert.Equal(t, "cat-123", category.Id())
				assert.Equal(t, "UpdatedCategory", category.Name())
				return expectedCategory, nil
			})

		// Act
		err := handler.UpdateCategory(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.UpdateCategoryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "cat-123", response.Category.Id)
		assert.Equal(t, "UpdatedCategory", response.Category.Name)
	})

	t.Run("異常系: IDが空", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":"UpdatedCategory"}`
		req := httptest.NewRequest(http.MethodPut, "/categories/", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.UpdateCategory(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("異常系: バリデーションエラー", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":""}`
		req := httptest.NewRequest(http.MethodPut, "/categories/cat-123", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("cat-123")

		// Act
		err := handler.UpdateCategory(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}

func TestCQRSServiceHandler_DeleteCategory(t *testing.T) {
	t.Run("正常系: カテゴリを削除できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/categories/cat-123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("cat-123")

		mockRepo.EXPECT().
			DeleteCategory(gomock.Any(), "cat-123").
			Return(nil)

		// Act
		err := handler.DeleteCategory(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("異常系: IDが空", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/categories/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.DeleteCategory(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("異常系: リポジトリエラー", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/categories/cat-123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("cat-123")

		mockRepo.EXPECT().
			DeleteCategory(gomock.Any(), "cat-123").
			Return(errors.New("database error"))

		// Act
		err := handler.DeleteCategory(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)
	})
}

func TestCQRSServiceHandler_CategoryById(t *testing.T) {
	t.Run("正常系: カテゴリを取得できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/categories/cat-123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("cat-123")

		expectedCategory := models.NewCategory("cat-123", "TestCategory")
		mockRepo.EXPECT().
			CategoryById(gomock.Any(), "cat-123").
			Return(expectedCategory, nil)

		// Act
		err := handler.CategoryById(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.CategoryByIdResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "cat-123", response.Category.Id)
		assert.Equal(t, "TestCategory", response.Category.Name)
	})

	t.Run("異常系: IDが空", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/categories/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.CategoryById(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}

func TestCQRSServiceHandler_CreateProduct(t *testing.T) {
	t.Run("正常系: 商品を作成できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":"TestProduct","price":1000,"category": {"id":"550e8400-e29b-41d4-a716-446655440000","name":"TestCategory"}}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// モックの設定
		category := models.NewCategory("550e8400-e29b-41d4-a716-446655440000", "TestCategory")
		expectedProduct := models.NewProduct("prod-123", "TestProduct", 1000, category)
		mockRepo.EXPECT().
			CreateProduct(gomock.Any(), "TestProduct", uint32(1000), gomock.Any()).
			Return(expectedProduct, nil)

		// Act
		err := handler.CreateProduct(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response dto.CreateProductResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "prod-123", response.Product.Id)
		assert.Equal(t, "TestProduct", response.Product.Name)
		assert.Equal(t, uint32(1000), response.Product.Price)
		assert.NotNil(t, response.Product.Category)
	})

	t.Run("異常系: バリデーションエラー（priceが0）", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":"TestProduct","price":0,"category_id":"550e8400-e29b-41d4-a716-446655440000"}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.CreateProduct(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("異常系: バリデーションエラー（category_idがUUID形式でない）", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":"TestProduct", "price":1000, "category": {"id":"invalid-uuid", "name":"TestCategory"}}`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.CreateProduct(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}

func TestCQRSServiceHandler_ProductList(t *testing.T) {
	t.Run("正常系: 商品一覧を取得できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// モックの設定
		category := models.NewCategory("cat-1", "Category1")
		expectedProducts := []*models.Product{
			models.NewProduct("prod-1", "Product1", 1000, category),
			models.NewProduct("prod-2", "Product2", 2000, category),
		}
		mockRepo.EXPECT().
			ProductList(gomock.Any()).
			Return(expectedProducts, nil)

		// Act
		err := handler.ProductList(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.ProductListResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response.Products, 2)
		assert.Equal(t, "prod-1", response.Products[0].Id)
		assert.Equal(t, uint32(1000), response.Products[0].Price)
	})

	t.Run("正常系: keywordパラメータで検索できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/products?keyword=test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// モックの設定（検索結果）
		category := models.NewCategory("cat-1", "Category1")
		expectedProducts := []*models.Product{
			models.NewProduct("prod-1", "TestProduct1", 1000, category),
			models.NewProduct("prod-2", "TestProduct2", 2000, category),
		}
		mockRepo.EXPECT().
			ProductByKeyword(gomock.Any(), "test").
			Return(expectedProducts, nil)

		// Act
		err := handler.ProductList(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.ProductByKeywordResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response.Products, 2)
		assert.Equal(t, "prod-1", response.Products[0].Id)
	})
}

func TestCQRSServiceHandler_ProductById(t *testing.T) {
	t.Run("正常系: 商品を取得できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/products/prod-123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("prod-123")

		category := models.NewCategory("cat-1", "Category1")
		expectedProduct := models.NewProduct("prod-123", "TestProduct", 1000, category)
		mockRepo.EXPECT().
			ProductById(gomock.Any(), "prod-123").
			Return(expectedProduct, nil)

		// Act
		err := handler.ProductById(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.ProductByIdResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "prod-123", response.Product.Id)
		assert.Equal(t, "TestProduct", response.Product.Name)
	})
}

func TestCQRSServiceHandler_UpdateProduct(t *testing.T) {
	t.Run("正常系: 商品を更新できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		e.Validator = server.NewCustomValidator()

		requestBody := `{"name":"UpdatedProduct","price":2000,"category":{"id":"550e8400-e29b-41d4-a716-446655440000","name":"TestCategory"}}`
		req := httptest.NewRequest(http.MethodPut, "/products/prod-123", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("prod-123")

		// モックの設定
		category := models.NewCategory("550e8400-e29b-41d4-a716-446655440000", "TestCategory")
		expectedProduct := models.NewProduct("prod-123", "UpdatedProduct", 2000, category)
		mockRepo.EXPECT().
			UpdateProduct(gomock.Any(), gomock.Any()).
			Return(expectedProduct, nil)

		// Act
		err := handler.UpdateProduct(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.UpdateProductResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "prod-123", response.Product.Id)
		assert.Equal(t, "UpdatedProduct", response.Product.Name)
		assert.Equal(t, uint32(2000), response.Product.Price)
	})
}

func TestCQRSServiceHandler_DeleteProduct(t *testing.T) {
	t.Run("正常系: 商品を削除できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/products/prod-123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("prod-123")

		mockRepo.EXPECT().
			DeleteProduct(gomock.Any(), "prod-123").
			Return(nil)

		// Act
		err := handler.DeleteProduct(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}

func TestCQRSServiceHandler_ProductByKeyword(t *testing.T) {
	t.Run("正常系: キーワードで商品を検索できる", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/products/search?keyword=test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// モックの設定
		category := models.NewCategory("cat-1", "Category1")
		expectedProducts := []*models.Product{
			models.NewProduct("prod-1", "TestProduct1", 1000, category),
			models.NewProduct("prod-2", "TestProduct2", 2000, category),
		}
		mockRepo.EXPECT().
			ProductByKeyword(gomock.Any(), "test").
			Return(expectedProducts, nil)

		// Act
		err := handler.ProductByKeyword(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.ProductByKeywordResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response.Products, 2)
	})

	t.Run("異常系: keywordが空", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockRepo := mock_repository.NewMockCQRSRepository(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		handler := server.NewCQRSServiceHandler(logger, mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/products/search", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.ProductByKeyword(c)

		// Assert
		require.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}
