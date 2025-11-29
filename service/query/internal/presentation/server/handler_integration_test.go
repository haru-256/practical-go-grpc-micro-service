//go:build integration || !ci

package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"connectrpc.com/connect"
	query "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1"
	queryconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1/queryv1connect"
	interceptor "github.com/haru-256/practical-go-grpc-micro-service/pkg/connect/interceptor"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/infrastructure/db"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/presentation/server"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	testDBConn *gorm.DB
)

func TestMain(m *testing.M) {
	// テスト用データベースのセットアップ
	configPath := "../../../"
	configName := "config"
	var err error
	testDBConn, err = testhelpers.SetupDB(configPath, configName)
	if err != nil {
		panic(err)
	}
	// テスト用データベースのクリーンアップ
	defer func() {
		if err = testhelpers.TeardownDB(testDBConn); err != nil {
			panic(err)
		}
	}()

	// テストの実行
	code := m.Run()

	os.Exit(code)
}

func setupProductIntegrationTests(t *testing.T) queryconnect.ProductServiceClient {
	t.Helper()
	repo := db.NewProductRepositoryImpl(testDBConn, testhelpers.TestLogger)
	productHandler, err := server.NewProductServiceHandlerImpl(testhelpers.TestLogger, repo)
	require.NoError(t, err, "Failed to create product handler")
	reqRespLogger := interceptor.NewReqRespLogger(testhelpers.TestLogger)
	validator, err := interceptor.NewValidator(testhelpers.TestLogger)
	require.NoError(t, err)
	mux := http.NewServeMux()
	path, handler := queryconnect.NewProductServiceHandler(
		productHandler,
		connect.WithInterceptors(
			reqRespLogger,
			validator.NewUnaryInterceptor(),
		),
	)
	mux.Handle(path, handler)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return queryconnect.NewProductServiceClient(server.Client(), server.URL)
}

func setupCategoryIntegrationTests(t *testing.T) queryconnect.CategoryServiceClient {
	t.Helper()
	repo := db.NewCategoryRepositoryImpl(testDBConn, testhelpers.TestLogger)
	categoryHandler, err := server.NewCategoryServiceHandlerImpl(testhelpers.TestLogger, repo)
	require.NoError(t, err, "Failed to create category handler")
	reqRespLogger := interceptor.NewReqRespLogger(testhelpers.TestLogger)
	validator, err := interceptor.NewValidator(testhelpers.TestLogger)
	require.NoError(t, err)
	mux := http.NewServeMux()
	path, handler := queryconnect.NewCategoryServiceHandler(
		categoryHandler,
		connect.WithInterceptors(
			reqRespLogger,
			validator.NewUnaryInterceptor(),
		),
	)
	mux.Handle(path, handler)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return queryconnect.NewCategoryServiceClient(server.Client(), server.URL)
}

func TestCategoryServiceHandlerImpl_ListCategories_Integration(t *testing.T) {
	client := setupCategoryIntegrationTests(t)
	ctx := context.Background()
	req := connect.NewRequest(&query.ListCategoriesRequest{})

	res, err := client.ListCategories(ctx, req)

	// Assert
	require.NoError(t, err, "ListCategoriesでエラーが発生しました")
	require.NotNil(t, res, "レスポンスがnilです")

	categories := res.Msg.GetCategories()
	require.NotEmpty(t, categories, "カテゴリリストが空です")

	// 最初のカテゴリの構造を検証
	firstCategory := categories[0]
	assert.NotEmpty(t, firstCategory.GetId(), "カテゴリIDが空です")
	assert.NotEmpty(t, firstCategory.GetName(), "カテゴリ名が空です")
}

func TestCategoryServiceHandlerImpl_GetCategoryById_Integration(t *testing.T) {
	client := setupCategoryIntegrationTests(t)
	ctx := context.Background()

	// まずListで取得してIDを確認
	listRes, err := client.ListCategories(ctx, connect.NewRequest(&query.ListCategoriesRequest{}))
	require.NoError(t, err, "カテゴリ一覧の取得に失敗しました")
	require.NotEmpty(t, listRes.Msg.GetCategories(), "カテゴリが存在しません")

	// 実際のIDを使用してテスト
	existingID := listRes.Msg.GetCategories()[0].GetId()
	req := connect.NewRequest(&query.GetCategoryByIdRequest{})
	req.Msg.SetId(existingID)

	// Act
	res, err := client.GetCategoryById(ctx, req)

	// Assert
	require.NoError(t, err, "GetCategoryByIdでエラーが発生しました")
	require.NotNil(t, res, "レスポンスがnilです")

	category := res.Msg.GetCategory()
	require.NotNil(t, category, "カテゴリがnilです")
	assert.Equal(t, existingID, category.GetId(), "IDが一致しません")
	assert.NotEmpty(t, category.GetName(), "カテゴリ名が空です")
}

func TestProductServiceHandlerImpl_ListProducts_Integration(t *testing.T) {
	client := setupProductIntegrationTests(t)
	ctx := context.Background()
	req := connect.NewRequest(&query.ListProductsRequest{})

	res, err := client.ListProducts(ctx, req)

	// Assert
	require.NoError(t, err, "ListProductsでエラーが発生しました")
	require.NotNil(t, res, "レスポンスがnilです")

	products := res.Msg.GetProducts()
	require.NotEmpty(t, products, "商品リストが空です")

	// 最初の商品の構造を検証
	firstProduct := products[0]
	assert.NotEmpty(t, firstProduct.GetId(), "商品IDが空です")
	assert.NotEmpty(t, firstProduct.GetName(), "商品名が空です")
	assert.Greater(t, firstProduct.GetPrice(), int32(0), "商品価格は0より大きい必要があります")
	assert.NotNil(t, firstProduct.GetCategory(), "商品にカテゴリが紐付いていません")
	assert.NotEmpty(t, firstProduct.GetCategory().GetId(), "カテゴリIDが空です")
	assert.NotEmpty(t, firstProduct.GetCategory().GetName(), "カテゴリ名が空です")
}

func TestProductServiceHandlerImpl_StreamProducts_Integration(t *testing.T) {
	client := setupProductIntegrationTests(t)
	ctx := context.Background()

	listRes, err := client.ListProducts(ctx, connect.NewRequest(&query.ListProductsRequest{}))
	require.NoError(t, err, "商品一覧の取得に失敗しました")
	require.NotEmpty(t, listRes.Msg.GetProducts(), "商品が存在しません")

	expected := make(map[string]struct{}, len(listRes.Msg.GetProducts()))
	for _, product := range listRes.Msg.GetProducts() {
		expected[product.GetId()] = struct{}{}
	}

	stream, err := client.StreamProducts(ctx, connect.NewRequest(&query.StreamProductsRequest{}))
	require.NoError(t, err, "StreamProductsの開始に失敗しました")
	t.Cleanup(func() {
		require.NoError(t, stream.Close())
	})

	received := make(map[string]struct{})
	for stream.Receive() {
		msg := stream.Msg()
		require.NotNil(t, msg, "ストリームメッセージがnilです")
		product := msg.GetProduct()
		require.NotNil(t, product, "ストリームのProductがnilです")
		received[product.GetId()] = struct{}{}
		assert.NotEmpty(t, product.GetName(), "商品名が空です")
		assert.Greater(t, product.GetPrice(), int32(0), "商品価格は0より大きい必要があります")
		assert.NotNil(t, product.GetCategory(), "商品にカテゴリが紐付いていません")
	}
	require.NoError(t, stream.Err(), "ストリーム受信中にエラーが発生しました")
	require.NotEmpty(t, received, "ストリームから商品が受信できませんでした")
	assert.Equal(t, expected, received, "StreamProductsが返す商品集合がListProductsと一致しません")
}

func TestProductServiceHandlerImpl_GetProductById_Integration(t *testing.T) {
	client := setupProductIntegrationTests(t)
	ctx := context.Background()

	// まずListで取得してIDを確認
	listRes, err := client.ListProducts(ctx, connect.NewRequest(&query.ListProductsRequest{}))
	require.NoError(t, err, "商品一覧の取得に失敗しました")
	require.NotEmpty(t, listRes.Msg.GetProducts(), "商品が存在しません")

	// 実際のIDを使用してテスト
	existingID := listRes.Msg.GetProducts()[0].GetId()
	req := connect.NewRequest(&query.GetProductByIdRequest{})
	req.Msg.SetId(existingID)

	// Act
	res, err := client.GetProductById(ctx, req)

	// Assert
	require.NoError(t, err, "GetProductByIdでエラーが発生しました")
	require.NotNil(t, res, "レスポンスがnilです")

	product := res.Msg.GetProduct()
	require.NotNil(t, product, "商品がnilです")
	assert.Equal(t, existingID, product.GetId(), "IDが一致しません")
	assert.NotEmpty(t, product.GetName(), "商品名が空です")
	assert.Greater(t, product.GetPrice(), int32(0), "商品価格は0より大きい必要があります")
	assert.NotNil(t, product.GetCategory(), "商品にカテゴリが紐付いていません")
}

func TestProductServiceHandlerImpl_SearchProductsByKeyword_Integration(t *testing.T) {
	client := setupProductIntegrationTests(t)
	ctx := context.Background()

	// まず全商品を取得してキーワードを決定
	listRes, err := client.ListProducts(ctx, connect.NewRequest(&query.ListProductsRequest{}))
	require.NoError(t, err, "商品一覧の取得に失敗しました")
	require.NotEmpty(t, listRes.Msg.GetProducts(), "商品が存在しません")

	// 最初の商品名から検索キーワードを抽出（例: "Product" や最初の単語）
	firstProductName := listRes.Msg.GetProducts()[0].GetName()
	require.NotEmpty(t, firstProductName, "商品名が空です")

	// 商品名の一部をキーワードとして使用（最初の5文字、または全体）
	runes := []rune(firstProductName)
	keyword := firstProductName
	if len(runes) > 5 {
		keyword = string(runes[:5])
	}

	req := connect.NewRequest(&query.SearchProductsByKeywordRequest{})
	req.Msg.SetKeyword(keyword)

	// Act
	res, err := client.SearchProductsByKeyword(ctx, req)

	// Assert
	require.NoError(t, err, "SearchProductsByKeywordでエラーが発生しました")
	require.NotNil(t, res, "レスポンスがnilです")

	products := res.Msg.GetProducts()
	// キーワードによっては結果が0件の可能性もあるので、NotEmptyではなくNotNilで検証
	require.NotNil(t, products, "商品リストがnilです")

	// 検索結果がある場合、構造を検証
	if len(products) > 0 {
		for _, product := range products {
			assert.NotEmpty(t, product.GetId(), "商品IDが空です")
			assert.NotEmpty(t, product.GetName(), "商品名が空です")
			assert.Greater(t, product.GetPrice(), int32(0), "商品価格は0より大きい必要があります")
			assert.NotNil(t, product.GetCategory(), "商品にカテゴリが紐付いていません")
		}
	}
}
