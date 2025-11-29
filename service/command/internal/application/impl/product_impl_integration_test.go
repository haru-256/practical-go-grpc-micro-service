//go:build integration || !ci

// Package impl provides integration tests for the application service layer.
// These tests verify the ProductService implementation including Add, Update, and Delete operations
// using a real database connection. Each test case is independent with automatic cleanup
// to ensure test isolation and prevent data pollution.
package impl

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/haru-256/practical-go-grpc-micro-service/pkg/errs"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/dto"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/application/service"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/categories"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/domain/models/products"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/infrastructure/sqlboiler/repository"
	"github.com/haru-256/practical-go-grpc-micro-service/service/command/internal/testhelpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProductService Integration Test", Ordered, func() {
	var (
		ps           service.ProductService
		cs           service.CategoryService
		tm           service.TransactionManager
		productRepo  products.ProductRepository
		categoryRepo categories.CategoryRepository
		ctx          context.Context
	)

	BeforeAll(func() {
		// データベース接続の初期化（categoryのsetupDatabaseを再利用）
		err := testhelpers.SetupDatabase("../../../", "config")
		Expect(err).NotTo(HaveOccurred())

		// テストではログ出力を破棄するloggerを使用
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		// サービスとリポジトリの初期化
		productRepo = repository.NewProductRepositoryImpl(logger)
		categoryRepo = repository.NewCategoryRepositoryImpl(logger)
		tm = repository.NewTransactionManagerImpl(logger)
		ps = NewProductServiceImpl(logger, productRepo, categoryRepo, tm)
		cs = NewCategoryServiceImpl(logger, categoryRepo, tm)
	})

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("Addメソッドの動作確認", func() {
		var (
			testProduct     *products.Product
			testCategory    *categories.Category
			testCategoryDTO *dto.CategoryDTO
		)

		BeforeEach(func() {
			var err error
			// 各テストで新しいユニークなカテゴリを作成
			CreateCategoryDTO := &dto.CreateCategoryDTO{
				Name: testhelpers.GenerateUniqueCategoryName(),
			}
			// カテゴリをデータベースに追加
			testCategoryDTO, err = cs.Add(ctx, CreateCategoryDTO)
			Expect(err).NotTo(HaveOccurred())
			testCategory, err = dto.CategoryFromDTO(testCategoryDTO)
			Expect(err).NotTo(HaveOccurred())

			// 各テストで新しいユニークな商品を作成
			name, err := products.NewProductName(testhelpers.GenerateUniqueProductName())
			Expect(err).NotTo(HaveOccurred(), "テスト用商品名の生成に失敗しました")
			price, err := products.NewProductPrice(1000)
			Expect(err).NotTo(HaveOccurred(), "テスト用商品価格の生成に失敗しました")
			testProduct, err = products.NewProduct(name, price, testCategory)
			Expect(err).NotTo(HaveOccurred(), "テスト用商品の生成に失敗しました")
		})

		It("新しい商品を追加できること", func() {
			createProductDTO := &dto.CreateProductDTO{
				Name: testProduct.Name().Value(),
				Category: &dto.CategoryDTO{
					Id:   testProduct.Category().Id().Value(),
					Name: testProduct.Category().Name().Value(),
				},
				Price: testProduct.Price().Value(),
			}

			result, err := ps.Add(ctx, createProductDTO)
			Expect(err).NotTo(HaveOccurred(), "商品の追加に失敗しました")
			Expect(result).NotTo(BeNil())
			Expect(result.Name).To(Equal(testProduct.Name().Value()))
			Expect(result.Price).To(Equal(uint32(testProduct.Price().Value())))
			DeferCleanup(testhelpers.CleanupProductCategory, tm, productRepo, categoryRepo, result, testCategoryDTO)

			// 追加された商品が存在することを確認
			exists, err := testhelpers.VerifyProductByName(ctx, tm, productRepo, testProduct.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "追加した商品がDBに存在しません")
		})

		It("既存の商品名で追加しようとするとエラーになること", func() {
			// 最初の追加
			createProductDTO := &dto.CreateProductDTO{
				Name: testProduct.Name().Value(),
				Category: &dto.CategoryDTO{
					Id:   testProduct.Category().Id().Value(),
					Name: testProduct.Category().Name().Value(),
				},
				Price: testProduct.Price().Value(),
			}
			testProductDTO, err := ps.Add(ctx, createProductDTO)
			Expect(err).NotTo(HaveOccurred())
			Expect(testProductDTO).NotTo(BeNil())
			DeferCleanup(testhelpers.CleanupProductCategory, tm, productRepo, categoryRepo, testProductDTO, testCategoryDTO)

			// 重複追加（同じ名前の新しい商品を作成）
			duplicateCreateProductDTO := &dto.CreateProductDTO{
				Name: testProduct.Name().Value(),
				Category: &dto.CategoryDTO{
					Id:   testProduct.Category().Id().Value(),
					Name: testProduct.Category().Name().Value(),
				},
				Price: testProduct.Price().Value(),
			}
			duplicateProductDTO, err := ps.Add(ctx, duplicateCreateProductDTO)
			Expect(err).To(HaveOccurred(), "重複する商品名で追加できてしまいました")
			Expect(duplicateProductDTO).To(BeNil())

			// エラーの詳細を検証
			appErr, ok := err.(*errs.ApplicationError)
			Expect(ok).To(BeTrue(), "返されたエラーがApplicationErrorではありません")
			Expect(appErr.Code).To(Equal("PRODUCT_ALREADY_EXISTS"), "エラーコードが期待値と異なります")
		})
	})

	Context("Updateメソッドの動作確認", func() {
		var (
			testProduct     *products.Product
			testCategory    *categories.Category
			testCategoryDTO *dto.CategoryDTO
			testProductDTO  *dto.ProductDTO
		)

		BeforeEach(func() {
			var err error
			// 各テストで新しいユニークなカテゴリを作成して追加
			createCategoryDTO := &dto.CreateCategoryDTO{
				Name: testhelpers.GenerateUniqueCategoryName(),
			}
			testCategoryDTO, err = cs.Add(ctx, createCategoryDTO)
			Expect(err).NotTo(HaveOccurred())
			testCategory, err = dto.CategoryFromDTO(testCategoryDTO)
			Expect(err).NotTo(HaveOccurred())

			// テスト用商品を作成して追加
			createProductDTO := &dto.CreateProductDTO{
				Name:  testhelpers.GenerateUniqueProductName(),
				Price: 1000,
				Category: &dto.CategoryDTO{
					Id:   testCategoryDTO.Id,
					Name: testCategoryDTO.Name,
				},
			}
			testProductDTO, err = ps.Add(ctx, createProductDTO)
			Expect(err).NotTo(HaveOccurred(), "テスト用商品の追加に失敗しました")
			testProduct, err = dto.ProductFromDTO(testProductDTO)
			Expect(err).NotTo(HaveOccurred())

			// クリーンアップ登録
			DeferCleanup(testhelpers.CleanupProductCategory, tm, productRepo, categoryRepo, testProductDTO, testCategoryDTO)
		})

		It("既存の商品を更新できること", func() {
			// 更新された商品を作成
			updatedName, err := products.NewProductName(testhelpers.GenerateUniqueProductName())
			Expect(err).NotTo(HaveOccurred(), "更新用商品名の生成に失敗しました")
			updatedPrice, err := products.NewProductPrice(2000)
			Expect(err).NotTo(HaveOccurred(), "更新用商品価格の生成に失敗しました")

			updateDTO := &dto.UpdateProductDTO{
				Id:         testProduct.Id().Value(),
				Name:       updatedName.Value(),
				CategoryId: testProduct.Category().Id().Value(),
				Price:      updatedPrice.Value(),
			}
			result, err := ps.Update(ctx, updateDTO)
			Expect(err).NotTo(HaveOccurred(), "商品の更新に失敗しました")
			Expect(result).NotTo(BeNil())
			Expect(result.Name).To(Equal(updatedName.Value()))
			Expect(result.Price).To(Equal(uint32(updatedPrice.Value())))

			// 更新された商品が存在することを確認
			exists, err := testhelpers.VerifyProductByName(ctx, tm, productRepo, updatedName)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue(), "更新した商品名がDBに存在しません")

			// 古い商品名が存在しないことを確認
			oldExists, err := testhelpers.VerifyProductByName(ctx, tm, productRepo, testProduct.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(oldExists).To(BeFalse(), "古い商品名がDBに残っています")
		})

		It("存在しない商品を更新しようとするとエラーになること", func() {
			// 存在しない商品を作成
			nonExistentName, err := products.NewProductName(testhelpers.GenerateUniqueProductName())
			Expect(err).NotTo(HaveOccurred())
			nonExistentPrice, err := products.NewProductPrice(3000)
			Expect(err).NotTo(HaveOccurred())
			nonExistentProduct, err := products.NewProduct(nonExistentName, nonExistentPrice, testCategory)
			Expect(err).NotTo(HaveOccurred())

			updateDTO := &dto.UpdateProductDTO{
				Id:         nonExistentProduct.Id().Value(),
				Name:       nonExistentProduct.Name().Value(),
				CategoryId: nonExistentProduct.Category().Id().Value(),
				Price:      nonExistentProduct.Price().Value(),
			}
			result, err := ps.Update(ctx, updateDTO)
			Expect(err).To(HaveOccurred(), "存在しない商品を更新できてしまいました")
			Expect(result).To(BeNil())

			// エラーの詳細を検証（リポジトリが返すエラー）
			crudErr, ok := err.(*errs.CRUDError)
			Expect(ok).To(BeTrue(), "返されたエラーがCRUDErrorではありません")
			Expect(crudErr.Code).To(Equal("NOT_FOUND"), "エラーコードが期待値と異なります")
		})
	})

	Context("Deleteメソッドの動作確認", func() {
		var (
			testCategory    *categories.Category
			testCategoryDTO *dto.CategoryDTO
			testProductDTO  *dto.ProductDTO
		)

		BeforeEach(func() {
			var err error
			// 各テストで新しいユニークなカテゴリを作成して追加
			createCategoryDTO := &dto.CreateCategoryDTO{
				Name: testhelpers.GenerateUniqueCategoryName(),
			}
			testCategoryDTO, err = cs.Add(ctx, createCategoryDTO)
			Expect(err).NotTo(HaveOccurred())
			testCategory, err = dto.CategoryFromDTO(testCategoryDTO)
			Expect(err).NotTo(HaveOccurred())

			// テスト用商品を作成して追加
			createProductDTO := &dto.CreateProductDTO{
				Name:  testhelpers.GenerateUniqueProductName(),
				Price: 1000,
				Category: &dto.CategoryDTO{
					Id:   testCategoryDTO.Id,
					Name: testCategoryDTO.Name,
				},
			}
			testProductDTO, err = ps.Add(ctx, createProductDTO)
			Expect(err).NotTo(HaveOccurred(), "テスト用商品の追加に失敗しました")

			// クリーンアップ登録
			DeferCleanup(testhelpers.CleanupProductCategory, tm, productRepo, categoryRepo, testProductDTO, testCategoryDTO)
		})

		It("既存の商品を削除できること", func() {
			deleteDTO := &dto.DeleteProductDTO{
				Id: testProductDTO.Id,
			}
			result, err := ps.Delete(ctx, deleteDTO)
			Expect(err).NotTo(HaveOccurred(), "商品の削除に失敗しました")
			Expect(result).NotTo(BeNil())
			Expect(result.Name).To(Equal(testProductDTO.Name))
			Expect(result.Price).To(Equal(uint32(testProductDTO.Price)))

			// 削除されたことを確認
			productID, err := products.NewProductId(result.Id)
			Expect(err).NotTo(HaveOccurred())
			exists, err := testhelpers.VerifyProductById(ctx, tm, productRepo, productID)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse(), "商品が正しく削除されていません")
		})

		It("存在しない商品を削除しようとするとエラーになること", func() {
			// 存在しない商品を作成
			nonExistentName, err := products.NewProductName(testhelpers.GenerateUniqueProductName())
			Expect(err).NotTo(HaveOccurred())
			nonExistentPrice, err := products.NewProductPrice(3000)
			Expect(err).NotTo(HaveOccurred())
			nonExistentProduct, err := products.NewProduct(nonExistentName, nonExistentPrice, testCategory)
			Expect(err).NotTo(HaveOccurred())

			deleteDTO := &dto.DeleteProductDTO{
				Id: nonExistentProduct.Id().Value(),
			}
			result, err := ps.Delete(ctx, deleteDTO)
			Expect(err).To(HaveOccurred(), "存在しない商品を削除できてしまいました")
			Expect(result).To(BeNil())

			// エラーの詳細を検証（リポジトリが返すエラー）
			var crudErr *errs.CRUDError
			if !errors.As(err, &crudErr) {
				// デバッグ用にエラーの型を出力
				GinkgoWriter.Printf("Error type: %T, Error: %v\n", err, err)
			}
			Expect(errors.As(err, &crudErr)).To(BeTrue(), "返されたエラーがCRUDErrorではありません")
			Expect(crudErr.Code).To(Equal("NOT_FOUND"), "エラーコードが期待値と異なります")
		})
	})
})
