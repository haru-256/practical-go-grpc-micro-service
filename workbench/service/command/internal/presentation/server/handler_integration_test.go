//go:build integration || !ci

package server_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"time"

	"connectrpc.com/connect"
	cmdconnect "github.com/haru-256/practical-go-grpc-micro-service/api/gen/go/command/v1/commandv1connect"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/impl"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/presentation/server"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/testhelpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	// testTimeout は各テストケースのタイムアウト
	testTimeout = 30 * time.Second
)

var _ = Describe("CategoryServiceHandler Integration Test", Label("IntegrationTests"), Ordered, func() {
	var (
		csh            *server.CategoryServiceHandlerImpl
		cs             service.CategoryService
		tm             service.TransactionManager
		repo           categories.CategoryRepository
		categoryClient cmdconnect.CategoryServiceClient
		ctx            context.Context
		cancel         context.CancelFunc
	)

	BeforeAll(func() {
		// データベース接続の初期化
		err := testhelpers.SetupDatabase("../../../", "config")
		Expect(err).NotTo(HaveOccurred())

		// テストではログ出力を破棄するloggerを使用
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		// サービスとリポジトリの初期化
		repo = repository.NewCategoryRepositoryImpl(logger)
		tm = repository.NewTransactionManagerImpl(logger)
		cs = impl.NewCategoryServiceImpl(logger, repo, tm)

		csh, err = server.NewCategoryServiceHandlerImpl(logger, cs)
		Expect(err).NotTo(HaveOccurred())

		validator, err := server.NewValidator(logger)
		Expect(err).NotTo(HaveOccurred())

		mux := http.NewServeMux()
		path, handler := cmdconnect.NewCategoryServiceHandler(
			csh,
			connect.WithInterceptors(validator.NewUnaryInterceptor()),
		)
		mux.Handle(path, handler)
		testServer := httptest.NewServer(mux)
		DeferCleanup(testServer.Close)

		categoryClient = cmdconnect.NewCategoryServiceClient(testServer.Client(), testServer.URL)
	})

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
		DeferCleanup(cancel)
	})

	Context("CreateCategory", func() {
		var testCategoryName *categories.CategoryName

		BeforeEach(func() {
			var err error
			// 各テストで新しいユニークなカテゴリを作成
			testCategoryName, err = categories.NewCategoryName(testhelpers.GenerateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred(), "テスト用カテゴリ名の生成に失敗しました")
		})

		It("カテゴリが正常に作成され、DBに永続化されること", func() {
			// Arrange
			req := testhelpers.CreateCategoryRequest(testCategoryName.Value())

			// Act
			resp, err := categoryClient.CreateCategory(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			testCategoryDTO := &dto.CategoryDTO{
				Id:   resp.Msg.GetCategory().GetId(),
				Name: resp.Msg.GetCategory().GetName(),
			}
			testCategory, err := dto.CategoryFromDTO(testCategoryDTO)
			Expect(err).NotTo(HaveOccurred(), "カテゴリDTOからドメインモデルへの変換に失敗しました")

			// テスト終了後のクリーンアップ
			DeferCleanup(testhelpers.CleanupCategory, tm, repo, testCategoryDTO)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Msg.GetCategory().GetName()).To(Equal(testCategoryName.Value()))
			Expect(resp.Msg.GetCategory().GetId()).NotTo(BeEmpty())
			Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())

			// DBに永続化されていることを確認
			exists, err := testhelpers.VerifyCategoryById(ctx, tm, repo, testCategory.Id())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "作成したカテゴリがDBに存在しません")
		})

		It("空の名前でバリデーションエラーを返すこと", func() {
			// Arrange
			req := testhelpers.CreateCategoryRequest("")

			// Act
			resp, err := categoryClient.CreateCategory(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
		})

		It("既に存在するカテゴリ名で作成するとエラーを返すこと", func() {
			// Arrange: 最初のカテゴリを作成
			req := testhelpers.CreateCategoryRequest(testCategoryName.Value())
			resp1, err := categoryClient.CreateCategory(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp1).NotTo(BeNil())
			testCategoryDTO := &dto.CategoryDTO{
				Id:   resp1.Msg.GetCategory().GetId(),
				Name: resp1.Msg.GetCategory().GetName(),
			}
			Expect(err).NotTo(HaveOccurred(), "カテゴリDTOからドメインモデルへの変換に失敗しました")

			// テスト終了後のクリーンアップ
			DeferCleanup(testhelpers.CleanupCategory, tm, repo, testCategoryDTO)

			// Act: 同じ名前で再度作成を試みる
			resp2, err := categoryClient.CreateCategory(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp2).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInternal))
		})
	})

	Context("UpdateCategory", func() {
		var testCategoryName *categories.CategoryName
		var createdCategoryId string

		BeforeEach(func() {
			var err error
			// 各テストで新しいユニークなカテゴリを作成
			testCategoryName, err = categories.NewCategoryName(testhelpers.GenerateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())

			// 更新用のカテゴリを事前に作成
			req := testhelpers.CreateCategoryRequest(testCategoryName.Value())
			resp, err := categoryClient.CreateCategory(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			// テスト終了後のクリーンアップ
			testCategoryDTO := &dto.CategoryDTO{
				Id:   resp.Msg.GetCategory().GetId(),
				Name: resp.Msg.GetCategory().GetName(),
			}
			DeferCleanup(testhelpers.CleanupCategory, tm, repo, testCategoryDTO)

			createdCategoryId = resp.Msg.GetCategory().GetId()
		})

		It("カテゴリが正常に更新され、DBに永続化されること", func() {
			// Arrange
			newName := testhelpers.GenerateUniqueCategoryName()
			req := testhelpers.CreateUpdateCategoryRequest(createdCategoryId, newName)

			// Act
			resp, err := categoryClient.UpdateCategory(ctx, req)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Msg.GetCategory().GetId()).To(Equal(createdCategoryId))
			Expect(resp.Msg.GetCategory().GetName()).To(Equal(newName))
			Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())

			// 更新後のカテゴリをDB検証
			updatedId, err := categories.NewCategoryId(createdCategoryId)
			Expect(err).NotTo(HaveOccurred())
			updatedName, err := categories.NewCategoryName(newName)
			Expect(err).NotTo(HaveOccurred())
			exists, err := testhelpers.VerifyCategoryById(ctx, tm, repo, updatedId)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "更新したカテゴリIDがDBに存在しません")
			exists, err = testhelpers.VerifyCategoryByName(ctx, tm, repo, updatedName)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "更新したカテゴリ名がDBに存在しません")
			exists, err = testhelpers.VerifyCategoryByName(ctx, tm, repo, testCategoryName)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse(), "古いカテゴリ名がDBに存在しています")
		})

		It("空のIDでバリデーションエラーを返すこと", func() {
			// Arrange
			req := testhelpers.CreateUpdateCategoryRequest("", "新しい名前")

			// Act
			resp, err := categoryClient.UpdateCategory(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
		})

		It("空の名前でバリデーションエラーを返すこと", func() {
			// Arrange
			req := testhelpers.CreateUpdateCategoryRequest(createdCategoryId, "")

			// Act
			resp, err := categoryClient.UpdateCategory(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
		})
	})

	Context("DeleteCategory", func() {
		var testCategory *categories.Category
		var createdCategoryId string

		BeforeEach(func() {
			// 各テストで新しいユニークなカテゴリを作成
			testCategoryName, err := categories.NewCategoryName(testhelpers.GenerateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())

			// 削除用のカテゴリを事前に作成
			req := testhelpers.CreateCategoryRequest(testCategoryName.Value())
			resp, err := categoryClient.CreateCategory(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			// テスト終了後のクリーンアップ
			testCategoryDTO := &dto.CategoryDTO{
				Id:   resp.Msg.GetCategory().GetId(),
				Name: resp.Msg.GetCategory().GetName(),
			}
			testCategory, err = dto.CategoryFromDTO(testCategoryDTO)
			Expect(err).NotTo(HaveOccurred())

			// テスト終了後のカテゴリクリーンアップ
			DeferCleanup(testhelpers.CleanupCategory, tm, repo, testCategoryDTO)

			createdCategoryId = resp.Msg.GetCategory().GetId()
		})

		It("カテゴリが正常に削除され、DBから削除されること", func() {
			// Arrange
			req := testhelpers.CreateDeleteCategoryRequest(createdCategoryId)

			// Act
			resp, err := categoryClient.DeleteCategory(ctx, req)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Msg.GetCategory().GetId()).To(Equal(createdCategoryId))
			Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())

			// DBから削除されていることを確認
			exists, _ := testhelpers.VerifyCategoryById(ctx, tm, repo, testCategory.Id())
			Expect(exists).To(BeFalse(), "削除したカテゴリがDBに存在します")
		})

		It("空のIDでバリデーションエラーを返すこと", func() {
			// Arrange
			req := testhelpers.CreateDeleteCategoryRequest("")

			// Act
			resp, err := categoryClient.DeleteCategory(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
		})
	})
})

var _ = Describe("ProductServiceHandler Integration Test", Label("IntegrationTests"), Ordered, func() {
	var (
		psh           *server.ProductServiceHandlerImpl
		ps            service.ProductService
		cs            service.CategoryService
		tm            service.TransactionManager
		productRepo   products.ProductRepository
		categoryRepo  categories.CategoryRepository
		productClient cmdconnect.ProductServiceClient
		ctx           context.Context
		cancel        context.CancelFunc
	)

	BeforeAll(func() {
		// データベース接続の初期化
		err := testhelpers.SetupDatabase("../../../", "config")
		Expect(err).NotTo(HaveOccurred())

		// テストではログ出力を破棄するloggerを使用
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		// サービスとリポジトリの初期化
		productRepo = repository.NewProductRepositoryImpl(logger)
		categoryRepo = repository.NewCategoryRepositoryImpl(logger)
		tm = repository.NewTransactionManagerImpl(logger)
		ps = impl.NewProductServiceImpl(logger, productRepo, categoryRepo, tm)
		cs = impl.NewCategoryServiceImpl(logger, categoryRepo, tm)

		psh, err = server.NewProductServiceHandlerImpl(logger, ps)
		Expect(err).NotTo(HaveOccurred())

		validator, err := server.NewValidator(logger)
		Expect(err).NotTo(HaveOccurred())

		mux := http.NewServeMux()
		path, handler := cmdconnect.NewProductServiceHandler(
			psh,
			connect.WithInterceptors(validator.NewUnaryInterceptor()),
		)
		mux.Handle(path, handler)
		testServer := httptest.NewServer(mux)
		DeferCleanup(testServer.Close)

		productClient = cmdconnect.NewProductServiceClient(testServer.Client(), testServer.URL)
	})

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
		DeferCleanup(cancel)
	})

	Context("CreateProduct", func() {
		var testProductName *products.ProductName
		var testCategory *categories.Category
		var testCategoryDTO *dto.CategoryDTO

		BeforeEach(func() {
			var err error
			// 各テストで新しいユニークな商品名を作成
			testProductName, err = products.NewProductName(testhelpers.GenerateUniqueProductName())
			Expect(err).NotTo(HaveOccurred(), "テスト用商品名の生成に失敗しました")

			// テスト用カテゴリを作成
			categoryName, err := categories.NewCategoryName(testhelpers.GenerateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())
			testCategory, err = categories.NewCategory(categoryName)
			Expect(err).NotTo(HaveOccurred())

			// カテゴリをDBに作成
			createCategoryDTO := &dto.CreateCategoryDTO{
				Name: testCategory.Name().Value(),
			}
			testCategoryDTO, err = cs.Add(ctx, createCategoryDTO)
			Expect(err).NotTo(HaveOccurred())

			// テスト終了後のカテゴリクリーンアップ
			DeferCleanup(testhelpers.CleanupCategory, tm, categoryRepo, testCategoryDTO)
		})

		It("商品が正常に作成され、DBに永続化されること", func() {
			// Arrange
			productPrice := uint32(1000)
			req := testhelpers.CreateProductRequest(
				testProductName.Value(),
				productPrice,
				testCategoryDTO.Id,
				testCategoryDTO.Name,
			)

			// Act
			resp, err := productClient.CreateProduct(ctx, req)
			testProductDTO := &dto.ProductDTO{
				Id:    resp.Msg.GetProduct().GetId(),
				Name:  resp.Msg.GetProduct().GetName(),
				Price: uint32(resp.Msg.GetProduct().GetPrice()),
				Category: &dto.CategoryDTO{
					Id:   resp.Msg.GetProduct().GetCategory().GetId(),
					Name: resp.Msg.GetProduct().GetCategory().GetName(),
				},
			}

			// テスト終了後のクリーンアップ
			DeferCleanup(testhelpers.CleanupProduct, tm, productRepo, testProductDTO)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Msg.GetProduct().GetName()).To(Equal(testProductName.Value()))
			Expect(resp.Msg.GetProduct().GetPrice()).To(Equal(int32(productPrice)))
			Expect(resp.Msg.GetProduct().GetCategory().GetId()).To(Equal(testCategoryDTO.Id))
			Expect(resp.Msg.GetProduct().GetCategory().GetName()).To(Equal(testCategoryDTO.Name))
			Expect(resp.Msg.GetProduct().GetId()).NotTo(BeEmpty())
			Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())

			// DBに永続化されていることを確認
			productId, err := products.NewProductId(testProductDTO.Id)
			Expect(err).NotTo(HaveOccurred())
			exists, err := testhelpers.VerifyProductById(ctx, tm, productRepo, productId)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "作成した商品がDBに存在しません")

			exists, err = testhelpers.VerifyProductByName(ctx, tm, productRepo, testProductName)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "作成した商品名がDBに存在しません")
		})

		It("空の商品名でバリデーションエラーを返すこと", func() {
			// Arrange
			req := testhelpers.CreateProductRequest(
				"",
				1000,
				testCategoryDTO.Id,
				testCategoryDTO.Name,
			)

			// Act
			resp, err := productClient.CreateProduct(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
		})

		It("価格が0でバリデーションエラーを返すこと", func() {
			// Arrange
			req := testhelpers.CreateProductRequest(
				testProductName.Value(),
				0,
				testCategoryDTO.Id,
				testCategoryDTO.Name,
			)

			// Act
			resp, err := productClient.CreateProduct(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
		})

		It("既に存在する商品名で作成するとエラーを返すこと", func() {
			// Arrange: 最初の商品を作成
			productPrice := uint32(1000)
			req := testhelpers.CreateProductRequest(
				testProductName.Value(),
				productPrice,
				testCategoryDTO.Id,
				testCategoryDTO.Name,
			)
			resp1, err := productClient.CreateProduct(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp1).NotTo(BeNil())

			testProductDTO := &dto.ProductDTO{
				Id:    resp1.Msg.GetProduct().GetId(),
				Name:  resp1.Msg.GetProduct().GetName(),
				Price: uint32(resp1.Msg.GetProduct().GetPrice()),
				Category: &dto.CategoryDTO{
					Id:   resp1.Msg.GetProduct().GetCategory().GetId(),
					Name: resp1.Msg.GetProduct().GetCategory().GetName(),
				},
			}

			// テスト終了後のクリーンアップ
			DeferCleanup(testhelpers.CleanupProduct, tm, productRepo, testProductDTO)

			// Act: 同じ名前で再度作成を試みる
			resp2, err := productClient.CreateProduct(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp2).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInternal))
		})
	})

	Context("UpdateProduct", func() {
		var testProduct *products.Product
		var testCategory *categories.Category
		var testCategoryDTO *dto.CategoryDTO
		var createdProductId string

		BeforeEach(func() {
			var err error
			// テスト用カテゴリを作成
			categoryName, err := categories.NewCategoryName(testhelpers.GenerateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())
			testCategory, err = categories.NewCategory(categoryName)
			Expect(err).NotTo(HaveOccurred())

			// カテゴリをDBに作成
			createCategoryDTO := &dto.CreateCategoryDTO{
				Name: testCategory.Name().Value(),
			}
			testCategoryDTO, err = cs.Add(ctx, createCategoryDTO)
			Expect(err).NotTo(HaveOccurred())

			// テスト終了後のカテゴリクリーンアップ
			DeferCleanup(testhelpers.CleanupCategory, tm, categoryRepo, testCategoryDTO)

			// 各テストで新しいユニークな商品を作成
			productName, err := products.NewProductName(testhelpers.GenerateUniqueProductName())
			Expect(err).NotTo(HaveOccurred())
			productPrice, err := products.NewProductPrice(1000)
			Expect(err).NotTo(HaveOccurred())
			testProduct, err = products.NewProduct(productName, productPrice, testCategory)
			Expect(err).NotTo(HaveOccurred())

			// 更新用の商品を事前に作成
			req := testhelpers.CreateProductRequest(
				testProduct.Name().Value(),
				testProduct.Price().Value(),
				testCategoryDTO.Id,
				testCategoryDTO.Name,
			)
			resp, err := productClient.CreateProduct(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			testProductDTO := &dto.ProductDTO{
				Id:    resp.Msg.GetProduct().GetId(),
				Name:  resp.Msg.GetProduct().GetName(),
				Price: uint32(resp.Msg.GetProduct().GetPrice()),
				Category: &dto.CategoryDTO{
					Id:   resp.Msg.GetProduct().GetCategory().GetId(),
					Name: resp.Msg.GetProduct().GetCategory().GetName(),
				},
			}

			// テスト終了後のクリーンアップ
			DeferCleanup(testhelpers.CleanupProduct, tm, productRepo, testProductDTO)

			createdProductId = resp.Msg.GetProduct().GetId()
		})

		It("商品が正常に更新され、DBに永続化されること", func() {
			// Arrange
			newName := testhelpers.GenerateUniqueProductName()
			newPrice := uint32(2000)
			req := testhelpers.UpdateProductRequest(
				createdProductId,
				newName,
				newPrice,
				testCategoryDTO.Id,
			)

			// Act
			resp, err := productClient.UpdateProduct(ctx, req)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Msg.GetProduct().GetId()).To(Equal(createdProductId))
			Expect(resp.Msg.GetProduct().GetName()).To(Equal(newName))
			Expect(resp.Msg.GetProduct().GetPrice()).To(Equal(int32(newPrice)))
			Expect(resp.Msg.GetProduct().GetCategory().GetId()).To(Equal(testCategoryDTO.Id))
			Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())

			// 更新後の商品をDB検証
			updatedId, err := products.NewProductId(createdProductId)
			Expect(err).NotTo(HaveOccurred())
			updatedName, err := products.NewProductName(newName)
			Expect(err).NotTo(HaveOccurred())

			exists, err := testhelpers.VerifyProductById(ctx, tm, productRepo, updatedId)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "更新した商品IDがDBに存在しません")

			exists, err = testhelpers.VerifyProductByName(ctx, tm, productRepo, updatedName)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "更新した商品名がDBに存在しません")

			exists, err = testhelpers.VerifyProductByName(ctx, tm, productRepo, testProduct.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse(), "古い商品名がDBに存在しています")
		})

		It("空のIDでバリデーションエラーを返すこと", func() {
			// Arrange
			req := testhelpers.UpdateProductRequest("", "新しい商品", 1000, testCategoryDTO.Id)

			// Act
			resp, err := productClient.UpdateProduct(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
		})

		It("空の商品名でバリデーションエラーを返すこと", func() {
			// Arrange
			req := testhelpers.UpdateProductRequest(createdProductId, "", 1000, testCategoryDTO.Id)

			// Act
			resp, err := productClient.UpdateProduct(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
		})

		It("価格が0でバリデーションエラーを返すこと", func() {
			// Arrange
			req := testhelpers.UpdateProductRequest(createdProductId, "商品名", 0, testCategoryDTO.Id)

			// Act
			resp, err := productClient.UpdateProduct(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
		})
	})

	Context("DeleteProduct", func() {
		var testProduct *products.Product
		var testCategory *categories.Category
		var testCategoryDTO *dto.CategoryDTO
		var createdProductId string

		BeforeEach(func() {
			var err error
			// テスト用カテゴリを作成
			categoryName, err := categories.NewCategoryName(testhelpers.GenerateUniqueCategoryName())
			Expect(err).NotTo(HaveOccurred())
			testCategory, err = categories.NewCategory(categoryName)
			Expect(err).NotTo(HaveOccurred())

			// カテゴリをDBに作成
			createCategoryDTO := &dto.CreateCategoryDTO{
				Name: testCategory.Name().Value(),
			}
			testCategoryDTO, err = cs.Add(ctx, createCategoryDTO)
			Expect(err).NotTo(HaveOccurred())

			// テスト終了後のカテゴリクリーンアップ
			DeferCleanup(testhelpers.CleanupCategory, tm, categoryRepo, testCategoryDTO)

			// 各テストで新しいユニークな商品を作成
			productName, err := products.NewProductName(testhelpers.GenerateUniqueProductName())
			Expect(err).NotTo(HaveOccurred())
			productPrice, err := products.NewProductPrice(1000)
			Expect(err).NotTo(HaveOccurred())
			testProduct, err = products.NewProduct(productName, productPrice, testCategory)
			Expect(err).NotTo(HaveOccurred())

			// 削除用の商品を事前に作成
			req := testhelpers.CreateProductRequest(
				testProduct.Name().Value(),
				testProduct.Price().Value(),
				testCategoryDTO.Id,
				testCategoryDTO.Name,
			)
			resp, err := productClient.CreateProduct(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			testProductDTO := &dto.ProductDTO{
				Id:    resp.Msg.GetProduct().GetId(),
				Name:  resp.Msg.GetProduct().GetName(),
				Price: uint32(resp.Msg.GetProduct().GetPrice()),
				Category: &dto.CategoryDTO{
					Id:   resp.Msg.GetProduct().GetCategory().GetId(),
					Name: resp.Msg.GetProduct().GetCategory().GetName(),
				},
			}

			// テスト終了後のクリーンアップ
			DeferCleanup(testhelpers.CleanupProduct, tm, productRepo, testProductDTO)

			createdProductId = resp.Msg.GetProduct().GetId()
		})

		It("商品が正常に削除され、DBから削除されること", func() {
			// Arrange
			req := testhelpers.DeleteProductRequest(createdProductId)

			// Act
			resp, err := productClient.DeleteProduct(ctx, req)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Msg.GetProduct().GetId()).To(Equal(createdProductId))
			Expect(resp.Msg.GetTimestamp()).NotTo(BeNil())

			// DBから削除されていることを確認
			productId, err := products.NewProductId(createdProductId)
			Expect(err).NotTo(HaveOccurred())
			exists, err := testhelpers.VerifyProductById(ctx, tm, productRepo, productId)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse(), "削除した商品がDBに存在します")
		})

		It("空のIDでバリデーションエラーを返すこと", func() {
			// Arrange
			req := testhelpers.DeleteProductRequest("")

			// Act
			resp, err := productClient.DeleteProduct(ctx, req)

			// Assert
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			var connectErr *connect.Error
			Expect(errors.As(err, &connectErr)).To(BeTrue())
			Expect(connectErr.Code()).To(Equal(connect.CodeInvalidArgument))
		})
	})
})
