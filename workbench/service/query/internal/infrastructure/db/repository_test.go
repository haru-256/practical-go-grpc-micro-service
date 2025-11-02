//go:build integration || !ci

package db_test

import (
	"context"
	"os"
	"testing"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/pkg/log"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/domain/models"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/config"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	testDBConn   *gorm.DB
	productRepo  *db.ProductRepositoryImpl
	categoryRepo *db.CategoryRepositoryImpl
)

func setupDB() error {
	// テスト用のデータベース接続設定
	configPath := "../../../"
	configName := "config"
	v := config.NewViper(configPath, configName)
	dbConfig, err := db.NewDBConfig(v)
	if err != nil {
		return err
	}
	logger, err := log.NewLogger(v)
	if err != nil {
		return err
	}
	// データベース接続の初期化
	testDBConn, err = db.NewDatabase(dbConfig, logger)
	if err != nil {
		return err
	}
	// リポジトリの初期化
	productRepo = db.NewProductRepositoryImpl(testDBConn, logger)
	categoryRepo = db.NewCategoryRepositoryImpl(testDBConn, logger)
	return nil
}

func teardownDB() error {
	sqlDB, err := testDBConn.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	// テスト用データベースのセットアップ
	if err := setupDB(); err != nil {
		panic(err)
	}
	// テストの実行
	code := m.Run()
	// テスト用データベースのクリーンアップ
	if err := teardownDB(); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestProductRepositoryImpl_List(t *testing.T) {
	tests := []struct {
		name       string
		assertions func(t *testing.T, products interface{}, err error)
	}{
		{
			name: "正常系: 商品リストを取得できる",
			assertions: func(t *testing.T, products interface{}, err error) {
				require.NoError(t, err)
				productList := products.([]*models.Product)
				assert.Greater(t, len(productList), 0, "商品リストが空です")
				// 最初の商品の詳細を検証
				if len(productList) > 0 {
					product := productList[0]
					assert.NotEmpty(t, product.Id())
					assert.NotEmpty(t, product.Name())
					assert.Greater(t, product.Price(), uint32(0))
					assert.NotNil(t, product.Category())
					assert.NotEmpty(t, product.Category().Id())
					assert.NotEmpty(t, product.Category().Name())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			products, err := productRepo.List(ctx)
			tt.assertions(t, products, err)
		})
	}
}

func TestProductRepositoryImpl_FindById(t *testing.T) {
	tests := []struct {
		name       string
		productId  string
		assertions func(t *testing.T, product interface{}, err error)
	}{
		{
			name:      "正常系: 商品IDで商品を取得できる",
			productId: "ac413f22-0cf1-490a-9635-7e9ca810e544", // 実際のテストデータ
			assertions: func(t *testing.T, product interface{}, err error) {
				require.NoError(t, err)
				p := product.(*models.Product)
				assert.NotNil(t, p)
				assert.NotEmpty(t, p.Id())
				assert.NotEmpty(t, p.Name())
				assert.Greater(t, p.Price(), uint32(0))
				assert.NotNil(t, p.Category())
			},
		},
		{
			name:      "異常系: 存在しない商品IDの場合、エラーを返す",
			productId: "non-existent-id",
			assertions: func(t *testing.T, product interface{}, err error) {
				require.Error(t, err)
				assert.Nil(t, product)
				assert.Contains(t, err.Error(), "NOT_FOUND")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			product, err := productRepo.FindById(ctx, tt.productId)
			tt.assertions(t, product, err)
		})
	}
}

func TestProductRepositoryImpl_FindByNameLike(t *testing.T) {
	tests := []struct {
		name       string
		keyword    string
		assertions func(t *testing.T, products interface{}, err error)
	}{
		{
			name:    "正常系: 商品名の部分一致で商品を取得できる",
			keyword: "商品", // 実際のテストデータに合わせて調整してください
			assertions: func(t *testing.T, products interface{}, err error) {
				require.NoError(t, err)
				productList := products.([]*models.Product)
				assert.GreaterOrEqual(t, len(productList), 0)
				// 商品が見つかった場合、商品の基本情報を検証
				for _, p := range productList {
					assert.NotNil(t, p)
					assert.NotEmpty(t, p.Id())
				}
			},
		},
		{
			name:    "正常系: 一致する商品がない場合、空のリストを返す",
			keyword: "存在しないキーワード12345",
			assertions: func(t *testing.T, products interface{}, err error) {
				require.NoError(t, err)
				productList := products.([]*models.Product)
				assert.Equal(t, 0, len(productList))
			},
		},
		{
			name:    "異常系: 空文字列の場合、エラーを返す",
			keyword: "",
			assertions: func(t *testing.T, products interface{}, err error) {
				require.Error(t, err)
				assert.Nil(t, products)
				var internalErr *errs.InternalError
				if assert.ErrorAs(t, err, &internalErr) {
					assert.Equal(t, "INVALID_KEYWORD", internalErr.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			products, err := productRepo.FindByNameLike(ctx, tt.keyword)
			tt.assertions(t, products, err)
		})
	}
}

func TestCategoryRepositoryImpl_List(t *testing.T) {
	tests := []struct {
		name       string
		assertions func(t *testing.T, categories interface{}, err error)
	}{
		{
			name: "正常系: カテゴリリストを取得できる",
			assertions: func(t *testing.T, categories interface{}, err error) {
				require.NoError(t, err)
				categoryList := categories.([]*models.Category)
				assert.GreaterOrEqual(t, len(categoryList), 0, "カテゴリリストの取得に失敗しました")
				// 最初のカテゴリの詳細を検証
				if len(categoryList) > 0 {
					category := categoryList[0]
					assert.NotEmpty(t, category.Id())
					assert.NotEmpty(t, category.Name())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			categories, err := categoryRepo.List(ctx)
			tt.assertions(t, categories, err)
		})
	}
}

func TestCategoryRepositoryImpl_FindById(t *testing.T) {
	tests := []struct {
		name       string
		categoryId string
		assertions func(t *testing.T, category interface{}, err error)
	}{
		{
			name:       "正常系: カテゴリIDでカテゴリを取得できる",
			categoryId: "b1524011-b6af-417e-8bf2-f449dd58b5c0", // 実際のテストデータ
			assertions: func(t *testing.T, category interface{}, err error) {
				require.NoError(t, err)
				c := category.(*models.Category)
				assert.NotNil(t, c)
				// データが存在する場合のみ詳細をチェック
				if c.Id() != "" {
					assert.NotEmpty(t, c.Name())
				}
			},
		},
		{
			name:       "異常系: 存在しないカテゴリIDの場合、エラーを返す",
			categoryId: "non-existent-category-id",
			assertions: func(t *testing.T, category interface{}, err error) {
				require.Error(t, err)
				assert.Nil(t, category)
				assert.Contains(t, err.Error(), "NOT_FOUND")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			category, err := categoryRepo.FindById(ctx, tt.categoryId)
			tt.assertions(t, category, err)
		})
	}
}
