package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	query "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1"
	queryconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/query/v1/queryv1connect"
	interceptor "github.com/haru-256/practical-go-grpc-micro-service/pkg/connect/interceptor"
	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/domain/models"
	mock_repository "github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/mock/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/query/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestProductServiceHandlerImpl_ListProducts はListProductsメソッドのテストです。
func TestProductServiceHandlerImpl_ListProducts(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*productHandlerSetup)
		wantErr      bool
		wantCode     connect.Code
		validateResp func(t *testing.T, resp *connect.Response[query.ListProductsResponse])
	}{
		{
			name: "正常系_商品リストが取得できる",
			setupMock: func(s *productHandlerSetup) {
				cat1 := models.NewCategory("cat1", "Electronics")
				cat2 := models.NewCategory("cat2", "Books")
				products := []*models.Product{
					models.NewProduct("prod1", "Product 1", 1000, cat1),
					models.NewProduct("prod2", "Product 2", 2000, cat2),
				}
				s.repo.EXPECT().List(gomock.Any()).Return(products, nil)
			},
			wantErr: false,
			validateResp: func(t *testing.T, resp *connect.Response[query.ListProductsResponse]) {
				require.NotNil(t, resp)
				products := resp.Msg.GetProducts()
				require.Len(t, products, 2)
				assert.Equal(t, "prod1", products[0].GetId())
				assert.Equal(t, "Product 1", products[0].GetName())
				assert.Equal(t, int32(1000), products[0].GetPrice())
			},
		},
		{
			name: "正常系_空のリストが取得できる",
			setupMock: func(s *productHandlerSetup) {
				s.repo.EXPECT().List(gomock.Any()).Return([]*models.Product{}, nil)
			},
			wantErr: false,
			validateResp: func(t *testing.T, resp *connect.Response[query.ListProductsResponse]) {
				require.NotNil(t, resp)
				assert.Empty(t, resp.Msg.GetProducts())
			},
		},
		{
			name: "異常系_リポジトリエラー",
			setupMock: func(s *productHandlerSetup) {
				s.repo.EXPECT().List(gomock.Any()).Return(nil, errs.NewInternalError("database", "database error"))
			},
			wantErr:  true,
			wantCode: connect.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setupProductHandler(t)
			defer s.cleanup()

			tt.setupMock(s)

			req := connect.NewRequest(&query.ListProductsRequest{})
			resp, err := s.client.ListProducts(s.ctx, req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.wantCode, connect.CodeOf(err))
			} else {
				require.NoError(t, err)
				if tt.validateResp != nil {
					tt.validateResp(t, resp)
				}
			}
		})
	}
}

// TestProductServiceHandlerImpl_StreamProducts はStreamProductsメソッドのテストです。
func TestProductServiceHandlerImpl_StreamProducts(t *testing.T) {
	readStream := func(t *testing.T, stream *connect.ServerStreamForClient[query.StreamProductsResponse]) []*query.StreamProductsResponse {
		t.Helper()
		defer func() {
			require.NoError(t, stream.Close())
		}()

		var messages []*query.StreamProductsResponse
		for stream.Receive() {
			msg := stream.Msg()
			require.NotNil(t, msg)
			messages = append(messages, msg)
		}
		require.NoError(t, stream.Err())

		return messages
	}

	tests := []struct {
		name            string
		setupMock       func(*productHandlerSetup)
		wantOpenErr     bool
		wantOpenErrCode connect.Code
		validate        func(t *testing.T, stream *connect.ServerStreamForClient[query.StreamProductsResponse])
		expectStreamErr bool
		streamErrCode   connect.Code
	}{
		{
			name: "正常系_複数の商品をストリーム受信できる",
			setupMock: func(s *productHandlerSetup) {
				cat1 := models.NewCategory("cat1", "Electronics")
				cat2 := models.NewCategory("cat2", "Books")
				products := []*models.Product{
					models.NewProduct("prod1", "Product 1", 1000, cat1),
					models.NewProduct("prod2", "Product 2", 2000, cat2),
				}
				s.repo.EXPECT().List(gomock.Any()).Return(products, nil)
			},
			validate: func(t *testing.T, stream *connect.ServerStreamForClient[query.StreamProductsResponse]) {
				messages := readStream(t, stream)
				require.Len(t, messages, 2)
				first := messages[0].GetProduct()
				second := messages[1].GetProduct()
				require.NotNil(t, first)
				require.NotNil(t, second)
				assert.Equal(t, "prod1", first.GetId())
				assert.Equal(t, "Product 1", first.GetName())
				assert.Equal(t, int32(1000), first.GetPrice())
				assert.Equal(t, "prod2", second.GetId())
				assert.Equal(t, "Product 2", second.GetName())
				assert.Equal(t, int32(2000), second.GetPrice())
			},
		},
		{
			name: "正常系_商品が存在しない場合は即終了する",
			setupMock: func(s *productHandlerSetup) {
				s.repo.EXPECT().List(gomock.Any()).Return([]*models.Product{}, nil)
			},
			validate: func(t *testing.T, stream *connect.ServerStreamForClient[query.StreamProductsResponse]) {
				messages := readStream(t, stream)
				assert.Len(t, messages, 0)
			},
		},
		{
			name: "異常系_リポジトリエラー",
			setupMock: func(s *productHandlerSetup) {
				s.repo.EXPECT().List(gomock.Any()).Return(nil, errs.NewInternalError("database", "database error"))
			},
			expectStreamErr: true,
			streamErrCode:   connect.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setupProductHandler(t)
			defer s.cleanup()

			if tt.setupMock != nil {
				tt.setupMock(s)
			}

			req := connect.NewRequest(&query.StreamProductsRequest{})
			stream, err := s.client.StreamProducts(s.ctx, req)

			if tt.wantOpenErr {
				require.Error(t, err)
				assert.Equal(t, tt.wantOpenErrCode, connect.CodeOf(err))
				return
			}

			require.NoError(t, err)
			if tt.expectStreamErr {
				received := stream.Receive()
				assert.False(t, received)
				require.Error(t, stream.Err())
				assert.Equal(t, tt.streamErrCode, connect.CodeOf(stream.Err()))
				require.NoError(t, stream.Close())
				return
			}

			if tt.validate != nil {
				tt.validate(t, stream)
			}
		})
	}
}

// TestProductServiceHandlerImpl_GetProductById はGetProductByIdメソッドのテストです。
func TestProductServiceHandlerImpl_GetProductById(t *testing.T) {
	tests := []struct {
		name         string
		productID    string
		setupMock    func(*productHandlerSetup)
		wantErr      bool
		wantCode     connect.Code
		validateResp func(t *testing.T, resp *connect.Response[query.GetProductByIdResponse])
	}{
		{
			name:      "正常系_商品が取得できる",
			productID: "prod1",
			setupMock: func(s *productHandlerSetup) {
				category := models.NewCategory("cat1", "Electronics")
				product := models.NewProduct("prod1", "Product 1", 1000, category)
				s.repo.EXPECT().FindById(gomock.Any(), "prod1").Return(product, nil)
			},
			wantErr: false,
			validateResp: func(t *testing.T, resp *connect.Response[query.GetProductByIdResponse]) {
				require.NotNil(t, resp)
				product := resp.Msg.GetProduct()
				require.NotNil(t, product)
				assert.Equal(t, "prod1", product.GetId())
				assert.Equal(t, "Product 1", product.GetName())
				assert.Equal(t, int32(1000), product.GetPrice())
			},
		},
		{
			name:      "異常系_商品が見つからない",
			productID: "nonexistent",
			setupMock: func(s *productHandlerSetup) {
				s.repo.EXPECT().FindById(gomock.Any(), "nonexistent").Return(nil, errs.NewCRUDError("NOT_FOUND", "product not found"))
			},
			wantErr:  true,
			wantCode: connect.CodeNotFound,
		},
		{
			name:      "異常系_バリデーションエラー_IDが空",
			productID: "",
			setupMock: func(s *productHandlerSetup) {
				// バリデーションエラーはリポジトリ呼び出し前に発生するため、モックは設定しない
			},
			wantErr:  true,
			wantCode: connect.CodeInvalidArgument,
		},
		{
			name:      "異常系_リポジトリエラー",
			productID: "prod1",
			setupMock: func(s *productHandlerSetup) {
				s.repo.EXPECT().FindById(gomock.Any(), "prod1").Return(nil, errs.NewInternalError("database", "database error"))
			},
			wantErr:  true,
			wantCode: connect.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setupProductHandler(t)
			defer s.cleanup()

			if tt.setupMock != nil {
				tt.setupMock(s)
			}

			req := connect.NewRequest(&query.GetProductByIdRequest{})
			req.Msg.SetId(tt.productID)
			resp, err := s.client.GetProductById(s.ctx, req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.wantCode, connect.CodeOf(err))
			} else {
				require.NoError(t, err)
				if tt.validateResp != nil {
					tt.validateResp(t, resp)
				}
			}
		})
	}
}

// TestProductServiceHandlerImpl_SearchProductsByKeyword はSearchProductsByKeywordメソッドのテストです。
func TestProductServiceHandlerImpl_SearchProductsByKeyword(t *testing.T) {
	tests := []struct {
		name         string
		keyword      string
		setupMock    func(*productHandlerSetup)
		wantErr      bool
		wantCode     connect.Code
		validateResp func(t *testing.T, resp *connect.Response[query.SearchProductsByKeywordResponse])
	}{
		{
			name:    "正常系_商品が検索できる",
			keyword: "Laptop",
			setupMock: func(s *productHandlerSetup) {
				category := models.NewCategory("cat1", "Electronics")
				products := []*models.Product{
					models.NewProduct("prod1", "Laptop Pro", 150000, category),
					models.NewProduct("prod2", "Laptop Air", 120000, category),
				}
				s.repo.EXPECT().FindByNameLike(gomock.Any(), "Laptop").Return(products, nil)
			},
			wantErr: false,
			validateResp: func(t *testing.T, resp *connect.Response[query.SearchProductsByKeywordResponse]) {
				require.NotNil(t, resp)
				require.Len(t, resp.Msg.GetProducts(), 2)
			},
		},
		{
			name:    "正常系_検索結果が空",
			keyword: "NonExistent",
			setupMock: func(s *productHandlerSetup) {
				s.repo.EXPECT().FindByNameLike(gomock.Any(), "NonExistent").Return([]*models.Product{}, nil)
			},
			wantErr: false,
			validateResp: func(t *testing.T, resp *connect.Response[query.SearchProductsByKeywordResponse]) {
				require.NotNil(t, resp)
				assert.Empty(t, resp.Msg.GetProducts())
			},
		},
		{
			name:    "異常系_バリデーションエラー_キーワードが空",
			keyword: "",
			setupMock: func(s *productHandlerSetup) {
				// バリデーションエラーはリポジトリ呼び出し前に発生
			},
			wantErr:  true,
			wantCode: connect.CodeInvalidArgument,
		},
		{
			name:    "異常系_リポジトリエラー",
			keyword: "Product",
			setupMock: func(s *productHandlerSetup) {
				s.repo.EXPECT().FindByNameLike(gomock.Any(), "Product").Return(nil, errs.NewInternalError("database", "database error"))
			},
			wantErr:  true,
			wantCode: connect.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setupProductHandler(t)
			defer s.cleanup()

			if tt.setupMock != nil {
				tt.setupMock(s)
			}

			req := connect.NewRequest(&query.SearchProductsByKeywordRequest{})
			req.Msg.SetKeyword(tt.keyword)
			resp, err := s.client.SearchProductsByKeyword(s.ctx, req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.wantCode, connect.CodeOf(err))
			} else {
				require.NoError(t, err)
				if tt.validateResp != nil {
					tt.validateResp(t, resp)
				}
			}
		})
	}
}

// TestCategoryServiceHandlerImpl_ListCategories はListCategoriesメソッドのテストです。
func TestCategoryServiceHandlerImpl_ListCategories(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*categoryHandlerSetup)
		wantErr      bool
		wantCode     connect.Code
		validateResp func(t *testing.T, resp *connect.Response[query.ListCategoriesResponse])
	}{
		{
			name: "正常系_カテゴリリストが取得できる",
			setupMock: func(s *categoryHandlerSetup) {
				categories := []*models.Category{
					models.NewCategory("cat1", "Electronics"),
					models.NewCategory("cat2", "Books"),
				}
				s.repo.EXPECT().List(gomock.Any()).Return(categories, nil)
			},
			wantErr: false,
			validateResp: func(t *testing.T, resp *connect.Response[query.ListCategoriesResponse]) {
				require.NotNil(t, resp)
				categories := resp.Msg.GetCategories()
				require.Len(t, categories, 2)
				assert.Equal(t, "cat1", categories[0].GetId())
				assert.Equal(t, "Electronics", categories[0].GetName())
			},
		},
		{
			name: "正常系_空のリストが取得できる",
			setupMock: func(s *categoryHandlerSetup) {
				s.repo.EXPECT().List(gomock.Any()).Return([]*models.Category{}, nil)
			},
			wantErr: false,
			validateResp: func(t *testing.T, resp *connect.Response[query.ListCategoriesResponse]) {
				require.NotNil(t, resp)
				assert.Empty(t, resp.Msg.GetCategories())
			},
		},
		{
			name: "異常系_リポジトリエラー",
			setupMock: func(s *categoryHandlerSetup) {
				s.repo.EXPECT().List(gomock.Any()).Return(nil, errs.NewInternalError("database", "database error"))
			},
			wantErr:  true,
			wantCode: connect.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setupCategoryHandler(t)
			defer s.cleanup()

			tt.setupMock(s)

			req := connect.NewRequest(&query.ListCategoriesRequest{})
			resp, err := s.client.ListCategories(s.ctx, req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.wantCode, connect.CodeOf(err))
			} else {
				require.NoError(t, err)
				if tt.validateResp != nil {
					tt.validateResp(t, resp)
				}
			}
		})
	}
}

// TestCategoryServiceHandlerImpl_GetCategoryById はGetCategoryByIdメソッドのテストです。
func TestCategoryServiceHandlerImpl_GetCategoryById(t *testing.T) {
	tests := []struct {
		name         string
		categoryID   string
		setupMock    func(*categoryHandlerSetup)
		wantErr      bool
		wantCode     connect.Code
		validateResp func(t *testing.T, resp *connect.Response[query.GetCategoryByIdResponse])
	}{
		{
			name:       "正常系_カテゴリが取得できる",
			categoryID: "cat1",
			setupMock: func(s *categoryHandlerSetup) {
				category := models.NewCategory("cat1", "Category 1")
				s.repo.EXPECT().FindById(gomock.Any(), "cat1").Return(category, nil)
			},
			wantErr: false,
			validateResp: func(t *testing.T, resp *connect.Response[query.GetCategoryByIdResponse]) {
				require.NotNil(t, resp)
				category := resp.Msg.GetCategory()
				require.NotNil(t, category)
				assert.Equal(t, "cat1", category.GetId())
				assert.Equal(t, "Category 1", category.GetName())
			},
		},
		{
			name:       "異常系_カテゴリが見つからない",
			categoryID: "nonexistent",
			setupMock: func(s *categoryHandlerSetup) {
				s.repo.EXPECT().FindById(gomock.Any(), "nonexistent").Return(nil, errs.NewCRUDError("NOT_FOUND", "category not found"))
			},
			wantErr:  true,
			wantCode: connect.CodeNotFound,
		},
		{
			name:       "異常系_バリデーションエラー_IDが空",
			categoryID: "",
			setupMock: func(s *categoryHandlerSetup) {
				// バリデーションエラーはリポジトリ呼び出し前に発生
			},
			wantErr:  true,
			wantCode: connect.CodeInvalidArgument,
		},
		{
			name:       "異常系_リポジトリエラー",
			categoryID: "cat1",
			setupMock: func(s *categoryHandlerSetup) {
				s.repo.EXPECT().FindById(gomock.Any(), "cat1").Return(nil, errs.NewInternalError("database", "database error"))
			},
			wantErr:  true,
			wantCode: connect.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setupCategoryHandler(t)
			defer s.cleanup()

			if tt.setupMock != nil {
				tt.setupMock(s)
			}

			req := connect.NewRequest(&query.GetCategoryByIdRequest{})
			req.Msg.SetId(tt.categoryID)
			resp, err := s.client.GetCategoryById(s.ctx, req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.wantCode, connect.CodeOf(err))
			} else {
				require.NoError(t, err)
				if tt.validateResp != nil {
					tt.validateResp(t, resp)
				}
			}
		})
	}
}

// productHandlerSetup はProductServiceHandlerのテスト用セットアップです。
type productHandlerSetup struct {
	ctx    context.Context
	ctrl   *gomock.Controller
	repo   *mock_repository.MockProductRepository
	client queryconnect.ProductServiceClient
	server *httptest.Server
}

// setupProductHandler はProductServiceHandlerのテストセットアップを作成します。
func setupProductHandler(t *testing.T) *productHandlerSetup {
	t.Helper()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockProductRepository(ctrl)

	handler, err := NewProductServiceHandlerImpl(testhelpers.TestLogger, repo)
	require.NoError(t, err)

	reqRespLogger := interceptor.NewReqRespLogger(testhelpers.TestLogger)
	validator, err := interceptor.NewValidator(testhelpers.TestLogger)
	require.NoError(t, err)

	mux := http.NewServeMux()
	path, handlerWithInterceptors := queryconnect.NewProductServiceHandler(
		handler,
		connect.WithInterceptors(
			reqRespLogger,
			validator.NewUnaryInterceptor(),
		),
	)
	mux.Handle(path, handlerWithInterceptors)
	testServer := httptest.NewServer(mux)

	client := queryconnect.NewProductServiceClient(testServer.Client(), testServer.URL)

	return &productHandlerSetup{
		ctx:    ctx,
		ctrl:   ctrl,
		repo:   repo,
		client: client,
		server: testServer,
	}
}

// cleanup はテストセットアップのクリーンアップを行います。
func (s *productHandlerSetup) cleanup() {
	s.server.Close()
	s.ctrl.Finish()
}

// categoryHandlerSetup はCategoryServiceHandlerのテスト用セットアップです。
type categoryHandlerSetup struct {
	ctx    context.Context
	ctrl   *gomock.Controller
	repo   *mock_repository.MockCategoryRepository
	client queryconnect.CategoryServiceClient
	server *httptest.Server
}

// setupCategoryHandler はCategoryServiceHandlerのテストセットアップを作成します。
func setupCategoryHandler(t *testing.T) *categoryHandlerSetup {
	t.Helper()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockCategoryRepository(ctrl)

	handler, err := NewCategoryServiceHandlerImpl(testhelpers.TestLogger, repo)
	require.NoError(t, err)

	reqRespLogger := interceptor.NewReqRespLogger(testhelpers.TestLogger)
	validator, err := interceptor.NewValidator(testhelpers.TestLogger)
	require.NoError(t, err)

	mux := http.NewServeMux()
	path, handlerWithInterceptors := queryconnect.NewCategoryServiceHandler(
		handler,
		connect.WithInterceptors(
			reqRespLogger,
			validator.NewUnaryInterceptor(),
		),
	)
	mux.Handle(path, handlerWithInterceptors)
	testServer := httptest.NewServer(mux)

	client := queryconnect.NewCategoryServiceClient(testServer.Client(), testServer.URL)

	return &categoryHandlerSetup{
		ctx:    ctx,
		ctrl:   ctrl,
		repo:   repo,
		client: client,
		server: testServer,
	}
}

// cleanup はテストセットアップのクリーンアップを行います。
func (s *categoryHandlerSetup) cleanup() {
	s.server.Close()
	s.ctrl.Finish()
}

func TestHandleError(t *testing.T) {
	t.Run("CRUDError NOT_FOUNDの場合CodeNotFoundを返す", func(t *testing.T) {
		err := errs.NewCRUDError("NOT_FOUND", "not found")
		grpcErr := handleError(err, "test operation")

		assert.Error(t, grpcErr)
		assert.Equal(t, connect.CodeNotFound, connect.CodeOf(grpcErr))
	})

	t.Run("CRUDError ALREADY_EXISTSの場合CodeAlreadyExistsを返す", func(t *testing.T) {
		err := errs.NewCRUDError("ALREADY_EXISTS", "already exists")
		grpcErr := handleError(err, "test operation")

		assert.Error(t, grpcErr)
		assert.Equal(t, connect.CodeAlreadyExists, connect.CodeOf(grpcErr))
	})

	t.Run("InternalErrorの場合CodeInternalを返す", func(t *testing.T) {
		err := errs.NewInternalError("INTERNAL", "internal error")
		grpcErr := handleError(err, "test operation")

		assert.Error(t, grpcErr)
		assert.Equal(t, connect.CodeInternal, connect.CodeOf(grpcErr))
	})

	t.Run("不明なエラーの場合CodeInternalを返す", func(t *testing.T) {
		err := errors.New("unknown error")
		grpcErr := handleError(err, "test operation")

		assert.Error(t, grpcErr)
		assert.Equal(t, connect.CodeInternal, connect.CodeOf(grpcErr))
	})
}

func TestToCategoryProto(t *testing.T) {
	t.Run("ドメインモデルからProtobufへの変換", func(t *testing.T) {
		category := models.NewCategory("cat1", "Category 1")
		proto := toCategoryProto(category)

		assert.Equal(t, "cat1", proto.GetId())
		assert.Equal(t, "Category 1", proto.GetName())
	})
}

func TestToCategoriesToProto(t *testing.T) {
	t.Run("ドメインモデルスライスからProtobufスライスへの変換", func(t *testing.T) {
		categories := []*models.Category{
			models.NewCategory("cat1", "Category 1"),
			models.NewCategory("cat2", "Category 2"),
		}
		protos := toCategoriesProto(categories)

		assert.Len(t, protos, 2)
		assert.Equal(t, "cat1", protos[0].GetId())
		assert.Equal(t, "Category 1", protos[0].GetName())
		assert.Equal(t, "cat2", protos[1].GetId())
		assert.Equal(t, "Category 2", protos[1].GetName())
	})

	t.Run("空のスライスの場合空のスライスを返す", func(t *testing.T) {
		categories := []*models.Category{}
		protos := toCategoriesProto(categories)

		assert.Empty(t, protos)
	})
}
