//go:build integration || !ci

package cqrs_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/domain/models"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/domain/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/infrastructure/config"
	"github.com/haru-256/practical-go-grpc-micro-service/service/client/internal/infrastructure/cqrs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	logger               *slog.Logger
	commandServiceClient *cqrs.CommandServiceClient
	queryServiceClient   *cqrs.QueryServiceClient
	pollingInterval      = 100 * time.Millisecond // polling interval for replication check
	replicationTimeout   = 5 * time.Second        // replication wait timeout
	streamTimeout        = 2 * time.Second        // timeout per StreamProducts attempt
)

func TestMain(m *testing.M) {
	logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	v := config.NewViper("../../../", "config")
	cfg, err := cqrs.NewCQRSServiceConfig(v)
	if err != nil {
		panic(err)
	}
	client := cqrs.NewClient(cfg)
	commandServiceClient = cqrs.NewCommandServiceClient(client, cfg)
	queryServiceClient = cqrs.NewQueryServiceClient(client, cfg)
	m.Run()
}

// Helper functions to reduce cognitive complexity

// waitForReplication waits for data to replicate to the query database
func waitForReplication(t *testing.T, assertFunc func() bool) {
	t.Helper()
	require.Eventually(t,
		assertFunc,
		replicationTimeout,
		pollingInterval,
		"指定時間内にレプリケーションされませんでした",
	)
}

// findProductByID checks if a product with the given ID exists in the slice
func findProductByID(products []*models.Product, targetID string) bool {
	for _, prod := range products {
		if prod.Id() == targetID {
			return true
		}
	}
	return false
}

// collectStreamProducts collects all products from a stream channel
func collectStreamProducts(t *testing.T, ctx context.Context, ch <-chan *repository.StreamProductsResult, timeout time.Duration) []*models.Product {
	t.Helper()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var products []*models.Product
	for {
		select {
		case <-ctxWithTimeout.Done():
			require.Fail(t, "StreamProducts did not complete before timeout")
			return nil
		case result, ok := <-ch:
			if !ok {
				return products
			}
			require.NoError(t, result.Err)
			products = append(products, result.Product)
		}
	}
}

// waitForProductReplication waits for a product to be available by ID
func waitForProductReplication(t *testing.T, ctx context.Context, repo *cqrs.CQRSRepositoryImpl, productID string, expectedName string, expectedPrice uint32) {
	t.Helper()
	waitForReplication(t, func() bool {
		product, err := repo.ProductById(ctx, productID)
		if err != nil {
			return false
		}
		require.NotNil(t, product)
		return assert.Equal(t, productID, product.Id()) &&
			assert.Equal(t, expectedName, product.Name()) &&
			assert.Equal(t, expectedPrice, product.Price())
	})
}

// assertProductInList waits for a product to appear in the product list
func assertProductInList(t *testing.T, ctx context.Context, repo *cqrs.CQRSRepositoryImpl, productID string) {
	t.Helper()
	waitForReplication(t, func() bool {
		products, err := repo.ProductList(ctx)
		require.NoError(t, err)
		assert.Greater(t, len(products), 0)
		return findProductByID(products, productID)
	})
}

// assertProductInStream waits for a product to appear in the streaming results
func assertProductInStream(t *testing.T, ctx context.Context, repo *cqrs.CQRSRepositoryImpl, productID string) {
	t.Helper()
	waitForReplication(t, func() bool {
		ch, err := repo.StreamProducts(ctx)
		require.NoError(t, err)
		require.NotNil(t, ch)

		products := collectStreamProducts(t, ctx, ch, streamTimeout)
		if !assert.Greater(t, len(products), 0) {
			return false
		}
		return findProductByID(products, productID)
	})
}

// assertProductInSearch waits for a product to appear in search results
func assertProductInSearch(t *testing.T, ctx context.Context, repo *cqrs.CQRSRepositoryImpl, keyword string, productID string) {
	t.Helper()
	waitForReplication(t, func() bool {
		products, err := repo.ProductByKeyword(ctx, keyword)
		require.NoError(t, err)
		return findProductByID(products, productID)
	})
}

func TestCQRSRepository_Category(t *testing.T) {
	ctx := t.Context()
	repo := cqrs.NewCQRSRepositoryImpl(
		commandServiceClient,
		queryServiceClient,
		logger,
	)

	var createdCategory *models.Category

	t.Run("カテゴリの作成", func(t *testing.T) {
		// カテゴリ名は1文字以上20文字以下
		categoryName := "TestCat" + uuid.New().String()[:8]
		created, err := repo.CreateCategory(ctx, categoryName)
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.Id())
		assert.Equal(t, categoryName, created.Name())
		createdCategory = created
	})

	t.Run("カテゴリIDで取得", func(t *testing.T) {
		require.NotNil(t, createdCategory, "カテゴリが作成されていません")
		// replication を待つためにEventuallyでリトライ
		require.Eventually(t,
			// この「条件関数」が true を返すまでリトライされる
			func() bool {
				// クエリ側（レプリカ側）のDBに問い合わせる
				category, err := repo.CategoryById(ctx, createdCategory.Id())

				// 1. そもそもエラー（例: not found）なら、まだレプリケーション中。
				//    false を返してリトライを継続する。
				if err != nil {
					return false
				}

				// 2. エラーがなく、期待通りの値か？
				require.NotNil(t, category)
				return assert.Equal(t, createdCategory.Id(), category.Id()) &&
					assert.Equal(t, createdCategory.Name(), category.Name())
			},
			// タイムアウト (最大待機時間)
			replicationTimeout,
			// ポーリング間隔
			pollingInterval,
			// 失敗した場合のメッセージ (オプション)
			"指定時間内にユーザーがレプリケーションされませんでした",
		)
	})

	t.Run("カテゴリの更新", func(t *testing.T) {
		require.NotNil(t, createdCategory, "カテゴリが作成されていません")

		// カテゴリ名は1文字以上20文字以下
		updatedName := "UpdateCat" + uuid.New().String()[:8]
		updatedCategory := models.NewCategory(createdCategory.Id(), updatedName)

		updated, err := repo.UpdateCategory(ctx, updatedCategory)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, createdCategory.Id(), updated.Id())
		assert.Equal(t, updatedName, updated.Name())
		createdCategory = updated
	})

	t.Run("カテゴリ一覧の取得", func(t *testing.T) {
		require.NotNil(t, createdCategory, "カテゴリが作成されていません")

		require.Eventually(t,
			func() bool {
				categories, err := repo.CategoryList(ctx)
				require.NoError(t, err)
				assert.Greater(t, len(categories), 0)
				// 作成したカテゴリが一覧に含まれることを確認
				found := false
				for _, cat := range categories {
					if cat.Id() == createdCategory.Id() {
						found = true
						break
					}
				}
				return assert.True(t, found, "作成したカテゴリが一覧に含まれていません")
			},
			replicationTimeout,
			pollingInterval,
			"指定時間内にユーザーがレプリケーションされませんでした",
		)
	})

	t.Run("カテゴリの削除", func(t *testing.T) {
		require.NotNil(t, createdCategory, "カテゴリが作成されていません")

		err := repo.DeleteCategory(ctx, createdCategory.Id())
		require.NoError(t, err)
	})
}

func TestCQRSRepository_Product(t *testing.T) {
	ctx := t.Context()
	repo := cqrs.NewCQRSRepositoryImpl(
		commandServiceClient,
		queryServiceClient,
		logger,
	)

	var testCategory *models.Category
	var createdProduct *models.Product

	// テスト用カテゴリの作成
	t.Run("テスト用カテゴリの作成", func(t *testing.T) {
		// カテゴリ名は1文字以上20文字以下
		categoryName := "ProdCat" + uuid.New().String()[:8]
		created, err := repo.CreateCategory(ctx, categoryName)
		require.NoError(t, err)
		require.NotNil(t, created)
		testCategory = created
	})

	t.Run("商品の作成", func(t *testing.T) {
		require.NotNil(t, testCategory, "テスト用カテゴリが作成されていません")

		// 商品名も適切な長さに
		productName := "TestProd" + uuid.New().String()[:8]
		productPrice := uint32(1000)
		created, err := repo.CreateProduct(ctx, productName, productPrice, testCategory)
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.Id())
		assert.Equal(t, productName, created.Name())
		assert.Equal(t, productPrice, created.Price())
		assert.Equal(t, testCategory.Id(), created.Category().Id())
		createdProduct = created
	})

	t.Run("商品IDで取得", func(t *testing.T) {
		require.NotNil(t, createdProduct, "商品が作成されていません")
		waitForProductReplication(t, ctx, repo, createdProduct.Id(), createdProduct.Name(), createdProduct.Price())
	})

	t.Run("商品の更新", func(t *testing.T) {
		require.NotNil(t, createdProduct, "商品が作成されていません")
		require.NotNil(t, testCategory, "テスト用カテゴリが作成されていません")

		// 商品名も適切な長さに
		updatedName := "UpdProd" + uuid.New().String()[:8]
		updatedPrice := uint32(2000)
		updatedProduct := models.NewProduct(createdProduct.Id(), updatedName, updatedPrice, testCategory)

		updated, err := repo.UpdateProduct(ctx, updatedProduct)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, createdProduct.Id(), updated.Id())
		assert.Equal(t, updatedName, updated.Name())
		assert.Equal(t, updatedPrice, updated.Price())
		createdProduct = updated
	})

	t.Run("商品一覧の取得", func(t *testing.T) {
		require.NotNil(t, createdProduct, "商品が作成されていません")
		assertProductInList(t, ctx, repo, createdProduct.Id())
	})

	t.Run("商品一覧の取得(Streaming)", func(t *testing.T) {
		require.NotNil(t, createdProduct, "商品が作成されていません")
		assertProductInStream(t, ctx, repo, createdProduct.Id())
	})

	t.Run("キーワードで商品検索", func(t *testing.T) {
		require.NotNil(t, createdProduct, "商品が作成されていません")
		keyword := "UpdProd"
		assertProductInSearch(t, ctx, repo, keyword, createdProduct.Id())
	})

	t.Run("商品の削除", func(t *testing.T) {
		require.NotNil(t, createdProduct, "商品が作成されていません")

		err := repo.DeleteProduct(ctx, createdProduct.Id())
		require.NoError(t, err)
	})

	// テスト用カテゴリのクリーンアップ
	t.Run("テスト用カテゴリの削除", func(t *testing.T) {
		if testCategory != nil {
			err := repo.DeleteCategory(ctx, testCategory.Id())
			require.NoError(t, err)
		}
	})
}
